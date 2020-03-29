package main

import (
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"
    "log"
    "io"
    "bufio"
    "os"
    "time"
    "fmt"
)

var (
    handle       *pcap.Handle
    device       string = "ens8"
    snapshot_len int32  = 1024
    promiscuous  bool   = false
    err          error
)


func constructWritingQueue( workers int ) chan packet_metadata {

    writingQueue := make(chan packet_metadata,100000)// 4* workers)
    return writingQueue
}

func constructZMapRoutine( workers int ) chan packet_metadata {


	//routine to read in from ZMap
	zmapIncoming := make(chan packet_metadata,100000)// 4*workers)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {

			//Read from ZMap
			input, err := reader.ReadString(byte('\n'))
			if err != nil && err == io.EOF {
				fmt.Println("ZMAP CLOSED")
                close(zmapIncoming)
				return
			}

            packet := convertToPacket( input )
            if packet == nil {
                continue
            }
			zmapIncoming <- *packet
		}

	}()

    return zmapIncoming
}

func constructPcapRoutine( workers int ) (chan packet_metadata, chan packet_metadata) {

	//routine to read in from pcap
	pcapIncoming := make(chan packet_metadata,1000000)//,4*workers )
	pcapdQueue := make(chan []byte,1000000)
	pcapQueue := make(chan packet_metadata,1000000)
	// Open device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, 1*time.Second)//pcap.BlockForever) //timeout
	if err != nil {
        panic(err)
		log.Fatal(err)
	}
	fmt.Println("1")
    for i := 0; i < workers/PARTITIONS; i ++ {
		go func(i int) {
			fmt.Println(i,"2")
            var eth layers.Ethernet
	        var ip4 layers.IPv4
	        var tcp layers.TCP
            var payload gopacket.Payload

            parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4,&tcp, &payload)
            pcapPacket := []gopacket.LayerType{}
			for {
				select {
				case data := <-pcapdQueue:
					if ip4.SrcIP.String() == "104.16.131.21" || ip4.SrcIP.String() == "104.16.131.20" {
						fmt.Println( "ah")
					}
					err := parser.DecodeLayers(data, &pcapPacket)
					if err == io.EOF {
						log.Println("Error:", err)
						return
					// packet does not have ipv4 or tcp
					} else if err != nil {
						continue
					}
					packet := ReadLayers( &ip4, &tcp )
					pcapIncoming <- *packet
				default:
					continue
				}
			}
        }(i)
    }
	go func() {
            for {
				select {
					case pcap:= <-pcapQueue:
						pcapIncoming <- pcap
					default:
						continue
				}
            }
    }()
    go func() {
		    defer handle.Close()

			for {
				// Use the handle as a packet source to process all packets
				data,_,_ := handle.ZeroCopyReadPacketData()
				pcapdQueue <- data
			}
			fmt.Println("stopped readingh dtya!!")
	}()

    return pcapIncoming, pcapQueue

}

func pollTimeoutRoutine( ipMeta * pState, timeoutQueue chan packet_metadata, workers int, timeout int ) (
    chan packet_metadata ) {

    TIMEOUT := time.Duration(timeout)*time.Second

	timeoutIncoming := make(chan packet_metadata,100000)//4*workers)
	timeoutQPass := make(chan packet_metadata,100000)//4*workers)
    //return from timeout when packet has expired
    go func() {
        for {
            select {
            case packet := <-timeoutQueue:

	            p, ok := ipMeta.find( &packet )
                //if no longer in map
	            if !ok {
                    //fmt.Println("no longer in map: " + string(packet.Saddr))
                    continue
                }
                //if timeout has reached, return packet.
                //else, check that the state has updated in the meanwhile
                //if not, put the packet back in timeoutQueue
                if (p.Counter > 1 && ( ((time.Now()).Sub( packet.Timestamp ) ) < TIMEOUT)) ||
                (((time.Now()).Sub( packet.Timestamp ) ) < 1*time.Second) {
                    //go func() { //must be its own routine to avoid deadlock
                        timeoutQPass <-packet
                        continue
                    //}()
                }else {
                    //fmt.Println("out of timeout")
                    //if state hasnt changed
                    if p.ExpectedRToLZR != packet.ExpectedRToLZR {
                        continue
                    } else {
                        //go func() { //must be its own routine to avoid deadlock
                            //fmt.Println("put into timeoutIncoming")
                            timeoutIncoming <-packet
                        //}()
                    }
                }
            case <-time.After(2 * time.Second):
                //fmt.Println("Something wrong with reading from timeoutQ")
                continue
            }
        }
    }()
    //dumb routine to avoid deadlock
    //pass to a passingQ
    go func() {
        for {
            select {
            case packet := <-timeoutQPass:
                timeoutQueue <- packet
            default:
                continue
           }
        }
    }()



    return timeoutIncoming

}

// TimeoutQueueStuff TODO:need to move
func constructTimeoutQueue( workers int ) chan packet_metadata {

    timeoutQueue := make(chan packet_metadata, 100000)
    return timeoutQueue
}



/*func constructIncomingChan() chan packet_metadata {

    incomingChan := make(chan packet_metadata)
    return incomingChan

}*/

