package lzr

import (
    "github.com/google/gopacket"
    //"github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"
    "log"
    "io"
    "bufio"
    "os"
    "time"
    "fmt"
	//"bytes"
)

var (
    handle       *pcap.Handle
    device       string = "ens8"
    snapshot_len int32  = 1024
    promiscuous  bool   = false
    err          error
)


func ConstructWritingQueue( workers int ) chan *packet_metadata {

    writingQueue := make(chan *packet_metadata,10000)// 4* workers)
    return writingQueue
}

func ConstructZMapRoutine( workers int ) chan *packet_metadata {


	//routine to read in from ZMap
	zmapIncoming := make(chan *packet_metadata,10000)// 4*workers)
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
			zmapIncoming <- packet
		}

	}()

    return zmapIncoming
}

func ConstructPcapRoutine( workers int ) (chan *packet_metadata, chan *packet_metadata) {

	//routine to read in from pcap
	pcapIncoming := make(chan *packet_metadata)//,10)//,4*workers )
	pcapdQueue := make(chan *gopacket.Packet)//,10)
	pcapQueue := make(chan *packet_metadata)//,10)
	// Open device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, pcap.BlockForever)//1*time.Second)
	if err != nil {
        panic(err)
		log.Fatal(err)
	}
	//set to filter out zmap syn packets (just syn) 
	err := handle.SetBPFFilter("tcp[tcpflags] != tcp-syn")
	if err != nil {
        panic(err)
		log.Fatal(err)
	}


    for i := 0; i < workers; i ++ {
		go func(i int) {
			for {
				select {
				case data := <-pcapdQueue:
					packet := convertToPacketM( data )
					if packet == nil {
						continue
					}
					pcapIncoming <- packet
				}
			}
        }(i)
    }
	go func() {
            for {
				select {
					case pcap:= <-pcapQueue:
						pcapIncoming <- pcap
				}
            }
    }()
    go func() {
		    defer handle.Close()
			packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
			for {
				pcapPacket, _ := packetSource.NextPacket()
				pcapdQueue <- &pcapPacket
			}
	}()

    return pcapIncoming, pcapQueue

}

func PollTimeoutRoutine( ipMeta * pState, timeoutQueue chan *packet_metadata, workers int, timeout int ) (
    chan *packet_metadata ) {

    TIMEOUT := time.Duration(timeout)*time.Second

	timeoutIncoming := make(chan *packet_metadata)//,1000)//4*workers)
	timeoutQPass := make(chan *packet_metadata)//,1000)//4*workers)
    //return from timeout when packet has expired
    go func() {
        for {
            select {
            case packet := <-timeoutQueue:

	            p, ok := ipMeta.find( packet )
                //if no longer in map
	            if !ok {
                    continue
                }
                //if timeout has reached, return packet.
                //else, check that the state has updated in the meanwhile
                //if not, put the packet back in timeoutQueue
                if !((p.Counter == 0 && ( ((time.Now()).Sub( packet.Timestamp ) ) > 1*time.Second)) ||
                (((time.Now()).Sub( packet.Timestamp ) ) > TIMEOUT)) {
                        timeoutQPass <-packet
                        continue
                }else {
                    //if state hasnt changed
                    if p.ExpectedRToLZR != packet.ExpectedRToLZR {
                        continue
                    } else {
                            timeoutIncoming <-packet
                    }
                }
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
           }
        }
    }()



    return timeoutIncoming

}

// TimeoutQueueStuff TODO:need to move
func ConstructTimeoutQueue( workers int ) chan *packet_metadata {

    timeoutQueue := make(chan *packet_metadata, 100000)
    return timeoutQueue
}

