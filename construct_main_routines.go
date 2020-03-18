package main

import (
    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
    "log"
    "io"
    "bufio"
    "os"
    "time"
    //"fmt"
)

var (
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

