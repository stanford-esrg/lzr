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


func constructWritingQueue( workers int ) chan packet_metadata {

    writingQueue := make(chan packet_metadata)//, workers)
    return writingQueue
}

func constructZMapRoutine( workers int ) chan packet_metadata {


	//routine to read in from ZMap
	zmapIncoming := make(chan packet_metadata ) //, workers)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {

			//Read from ZMap
			input, err := reader.ReadString(byte('\n'))
			if err != nil && err == io.EOF {
                close(zmapIncoming)
				return
			}
            packet := convertToPacket( input )
            if packet == nil {
                return
            }
			zmapIncoming <- *packet
		}

	}()

    return zmapIncoming
}

func constructPcapRoutine( workers int ) chan packet_metadata {

	//routine to read in from pcap
	pcapIncoming := make(chan packet_metadata) //, workers)
	go func() {
		// Open device
		handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, 0) //timeout
		if err != nil {
            //panic(err)
			log.Fatal(err)
		}
		defer handle.Close()
		// Use the handle as a packet source to process all packets
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

		for {

			//Read from pcap
			pcapPacket, err := packetSource.NextPacket()
			if err == io.EOF {
				return
			} else if err != nil {
				log.Println("Error:", err)
				continue
			}
            packet := convertToPacketM( pcapPacket )
            if packet == nil {
                continue
            }
			pcapIncoming <- *packet
		}

	}()

    return pcapIncoming

}



func pollTimeoutRoutine( ipMeta * pState, timeoutQueue chan packet_metadata, workers int, timeout int ) chan packet_metadata {

    TIMEOUT := time.Duration(timeout)*time.Second

	timeoutIncoming := make(chan packet_metadata)//, workers)
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
                //fmt.Println("out of timeout")
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
                        //fmt.Println("put into timeoutIncoming")
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

    timeoutQueue := make(chan packet_metadata)//, workers)
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

    writingQueue := constructWritingQueue( options.Workers )
    zmapIncoming := constructZMapRoutine( options.Workers )
    pcapIncoming := constructPcapRoutine( options.Workers )
    timeoutQueue := constructTimeoutQueue( options.Workers )
    timeoutIncoming := pollTimeoutRoutine( &ipMeta,timeoutQueue, options.Workers, options.Timeout )
    timeoutEmpty := true

    // record to file
    go func() {
        for {
            select {
                case input := <-writingQueue:
                    f.record( input )
                default:
                    continue
                }
        }
    }()


	//read from both zmap and pcap
	for {
		select {
			case input, ok := <-zmapIncoming:
                //ExitCondition:
                //done reading from zmap (channel closed)
                //TODO: timeoutQ is empty ??????
                //eventually check: all locks are released (no more jobs runnign)
                if !ok && timeoutEmpty {
                    time.Sleep( 1 * time.Second )
                    return
                }
                if err := sem.Acquire(ctx, 1); err != nil {
                    continue
                }
                go func() {
                    defer sem.Release(1)
                    // not checking if another thread is processing since we are 
                    // assuming that repeats are being filtered at zmap 
                    // and thus IPs cannot be in ipMeta b4 zmap adds to it
				    ackZMap( input, &ipMeta, &timeoutQueue, &writingQueue, f )
                }()
			case input := <-pcapIncoming:

                if err := sem.Acquire(ctx, 1); err != nil {
                    continue
                }
                go func() {
                    defer sem.Release(1)
                    inMap, processing := ipMeta.isProcessing( &input ) 
                    //if another thread is processing, put input back
                    if processing {
                        pcapIncoming <- input
                        return
                    }
                    //if not in map, return
                    if !inMap {
                        return
                    }
                    ipMeta.startProcessing( &input )
				    handlePcap( input, &ipMeta, &timeoutQueue, &writingQueue, f )
                    ipMeta.finishProcessing( &input )
                }()
            case input, _ := <-timeoutIncoming:
                if err := sem.Acquire(ctx, 1); err != nil {
                    continue
                }
                go func() {
                    defer sem.Release(1)
                    inMap, processing := ipMeta.isProcessing( &input ) 
                    //if another thread is processing, put input back
                    if processing {
                        pcapIncoming <- input
                        return
                    }
                    //if not in map, return
                    if !inMap {
                        return
                    }
                    ipMeta.startProcessing( &input )
                    handleTimeout( input, &ipMeta, &timeoutQueue, &writingQueue, f )
                    ipMeta.finishProcessing( &input )
                }()

			default:
                continue
		}
	}
} //end of main
