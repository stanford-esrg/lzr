package main

import (
    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
    "log"
    "io"
    "bufio"
    "os"
    "time"
    "context"
    "golang.org/x/sync/semaphore"
    //"fmt"
)


//TODO: move vars to appropriate places
var (
    device       string = "ens8"
    snapshot_len int32  = 1024
    promiscuous  bool   = false
    err          error
    handle       *pcap.Handle
)

func constructZMapRoutine( workers int ) chan string {


	//routine to read in from ZMap
	zmapIncoming := make(chan string, workers)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {

			//Read from ZMap
			input, err := reader.ReadString(byte('\n'))
			if err != nil && err == io.EOF {
				return
			}
			zmapIncoming <- input
		}

	}()

    return zmapIncoming
}

func constructPcapRoutine( workers int ) chan gopacket.Packet {

	//routine to read in from pcap
	pcapIncoming := make(chan gopacket.Packet, workers)
	go func() {
		// Open device
		handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, 0) //timeout
		if err != nil {
			log.Fatal(err)
		}
		defer handle.Close()
		// Use the handle as a packet source to process all packets
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

		for {

			//Read from pcap
			packet, err := packetSource.NextPacket()
			if err == io.EOF {
				return
			} else if err != nil {
				log.Println("Error:", err)
				continue
			}
			pcapIncoming <- packet
		}

	}()

    return pcapIncoming

}



func pollTimeoutRoutine( ipMeta * pState, timeoutQueue chan packet_metadata, workers int, timeout int ) chan packet_metadata {

    TIMEOUT := time.Duration(timeout)*time.Second

	timeoutIncoming := make(chan packet_metadata, workers)
    //timeoutReQ := make(chan packet_metadata) //to avoid deadlock need 
    //return from timeout when packet has expired
    go func() {
        for {
            packet := <-timeoutQueue
            //if timeout has reached, return packet.
            //else, check that the state has updated in the meanwhile
            //if not, put the packet back in timeoutQueue
            if ( ((time.Now()).Sub( packet.Timestamp ) ) < TIMEOUT) {
                go func() { //must be its own routine to avoid deadlock
                    timeoutQueue <-packet
                }()
            } else {
	            p, ok := ipMeta.find( &packet )
                //if no longer in map
	            if !ok {
                    //fmt.Println("no longer in map: " + string(packet.Saddr))
                    continue
                }
                //if state hasnt changed
                if p.ExpectedR != packet.ExpectedR {
                    continue
                } else {
                    go func() { //must be its own routine to avoid deadlock
                        timeoutIncoming <-packet
                    }()
                }
            }
        }
    }()
    return timeoutIncoming

}

// TimeoutQueueStuff TODO:need to move
func constructTimeoutQueue( workers int ) chan packet_metadata {

    timeoutQueue := make(chan packet_metadata, workers)
    return timeoutQueue
}



/*func constructIncomingChan() chan packet_metadata {

    incomingChan := make(chan packet_metadata)
    return incomingChan

}*/

func main() {

    ctx := context.TODO()

    //read in config 
    options := parse()

	//initalize
	ipMeta := constructPacketStateMap()
    f := initFile( options.Filename )
    sem := semaphore.NewWeighted( int64(options.Workers) )

    zmapIncoming := constructZMapRoutine( options.Workers )
    pcapIncoming := constructPcapRoutine( options.Workers )
    timeoutQueue := constructTimeoutQueue( options.Workers )
    timeoutIncoming := pollTimeoutRoutine( &ipMeta,timeoutQueue, options.Workers, options.Timeout )

	//read from both zmap and pcap
	for {
		select {
			case input := <-zmapIncoming:
                if err := sem.Acquire(ctx, 1); err != nil {
                    continue
                }
                go func() {
                    defer sem.Release(1)
				    ackZMap( input, &ipMeta, &timeoutQueue )
                }()
			case input := <-pcapIncoming:
                if err := sem.Acquire(ctx, 1); err != nil {
                    continue
                }
                go func() {
                    defer sem.Release(1)
				    handlePcap( input, &ipMeta, &timeoutQueue, f )
                }()
            case input := <-timeoutIncoming:
                if err := sem.Acquire(ctx, 1); err != nil {
                    continue
                }
                go func() {
                    defer sem.Release(1)
                    handleTimeout( input, &ipMeta, &timeoutQueue, f )
                }()
			default:
                continue
		}
	}
    /*
    for i := 0; i < options.Workers; i++ {
		go func( i int ) {
	        for {
	            //read from both zmap and pcap and timeout
		        select {
			        case input := <-zmapIncoming:
				        ackZMap( input, &ipMeta, &timeoutQueue )
			        case input := <-pcapIncoming:
				        handlePcap( input, &ipMeta, &timeoutQueue, f )
                    case input := <-timeoutIncoming:
                        handleTimeout( input, &ipMeta, &timeoutQueue, f )
			        default:
                        continue
                }
            }
		}( i )
	}

    //temp solution to wait forever
    select{}
    */
} //end of main
