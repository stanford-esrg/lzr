package main

import (
    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
    "log"
    "io"
    "bufio"
    "os"
    "time"
    "fmt"
)

//TODO: move vars to appropriate places
var (
    device       string = "ens8"
    snapshot_len int32  = 1024
    promiscuous  bool   = false
    err          error
    timeout      time.Duration = 5 * time.Second
    handle       *pcap.Handle
    buffer       gopacket.SerializeBuffer
)

func constructZMapRoutine() chan string {


	//routine to read in from ZMap
	zmapIncoming := make(chan string)
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

func constructPcapRoutine() chan gopacket.Packet {

	//routine to read in from pcap
	pcapIncoming := make(chan gopacket.Packet)
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



func pollTimeoutRoutine( ipMeta * pState, timeoutQueue chan packet_metadata ) chan packet_metadata {

    TIMEOUT := 2*time.Second

	timeoutIncoming := make(chan packet_metadata)
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
	            p, ok := (*ipMeta).IPmap[packet.Saddr]
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
func constructTimeoutQueue() chan packet_metadata {

    timeoutQueue := make(chan packet_metadata)
    return timeoutQueue
}



func constructIncomingChan() chan packet_metadata {

    incomingChan := make(chan packet_metadata)
    return incomingChan

}

func main() {

    //read in config 
    options := parse()

	//initalize
	ipMeta := constructPacketStateMap()
    timeoutQueue := constructTimeoutQueue()
    f := initFile( options.Filename )

    //TODO: merge all channels into one?
    //incomingChan := constructIncomingChan()
    zmapIncoming := constructZMapRoutine()
    pcapIncoming := constructPcapRoutine()
    timeoutIncoming := pollTimeoutRoutine( &ipMeta,timeoutQueue )

	//read from both zmap and pcap
	for {
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

} //end of main
