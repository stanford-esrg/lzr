/*
Copyright 2020 The Board of Trustees of The Leland Stanford Junior University

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package lzr

import (
	"log"

    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"

    //"io"
    "bufio"
	"fmt"
    "os"
	"strings"
    "time"
)

var (
    handle       *pcap.Handle
    snapshot_len int32  = 1024
    promiscuous  bool   = false
    err          error
    source_mac   string
    dest_mac     string
	QUEUE_SIZE   int32 = 200000000
)

func InitParams() {

    source_mac = getSourceMacAddr()
    dest_mac =  getHostMacAddr()

}


func ConstructWritingQueue( workers int ) chan *packet_metadata {

    writingQueue := make(chan *packet_metadata, QUEUE_SIZE)
    return writingQueue
}


func ConstructIncomingRoutine( workers int ) chan *packet_metadata {
	incoming := make(chan *packet_metadata, QUEUE_SIZE)
	go func() {
		defer close(incoming)
		scanner := bufio.NewScanner(os.Stdin)
		var destIP, destPort string

		var ticker *time.Ticker
		if !ReadZMap() && !DryRun() {
			ticker = time.NewTicker(time.Second)
			defer ticker.Stop()
		}

		buffer := make([]*packet_metadata, 0, *rate)

		for scanner.Scan() {
			input := scanner.Text()
			var packet *packet_metadata

			if ReadZMap () {
				packet = convertFromZMapToPacket( input )
			} else if DryRun() {
				// Process ZMap dry-run output
				if strings.Contains(input, "daddr:") {
					parts := strings.Fields(input)
					for i, part := range parts {
						if part == "daddr:" && i+1 < len(parts) {
							destIP = parts[i+1]
							break
						}
					}
				} else if strings.Contains(input, "dest:") {
					parts := strings.Fields(input)
					for i, part := range parts {
						if part == "dest:" && i+1 < len(parts) {
							destPort = strings.TrimSuffix(parts[i+1], "|")
							break
						}
					}

					if destIP != "" && destPort != "" {
						ipPortString := fmt.Sprintf("%s:%s", destIP, destPort)
						packet = convertFromInputListToPacket(ipPortString)
						destIP, destPort = "", ""
					}
				}
			} else {
				
				packet = convertFromInputListToPacket( input )
			}

			if packet != nil {
				buffer = append(buffer, packet)
			}
			if ticker != nil && len(buffer) >= *rate {
				<-ticker.C
				for _, pkt := range buffer {
					incoming <- pkt
				}
				// Clear buffer after dispatch
				buffer = buffer[:0]
			}
		}

		// Dispatch any remaining packets (if any)
		if len(buffer) > 0 {
			if ticker != nil {
				<-ticker.C
			}
			for _, pkt := range buffer {
				incoming <- pkt
			}
		}
		if ticker != nil {
			<-ticker.C
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		}
		fmt.Fprintln(os.Stderr, "Finished Reading Input")
	}()

    return incoming
}

func ConstructPcapRoutine( workers int, useIPv6 bool ) chan *packet_metadata {

	//routine to read in from pcap
	pcapIncoming := make(chan *packet_metadata, QUEUE_SIZE)
	pcapdQueue := make(chan *gopacket.Packet, QUEUE_SIZE)
	// Open device
	handle, err = pcap.OpenLive(getDevice(), snapshot_len, promiscuous, pcap.BlockForever)//1*time.Second)
	if err != nil {
        panic(err)
		log.Fatal(err)
	}
	//set to filter out zmap syn packets (just syn) 
	filter := "tcp[tcpflags] != tcp-syn"
	if useIPv6 == true {
		filter = "(ip6 proto 6 && (ip6[53] & 4 != 0 || ip6[53] != 2))"
	}
	err := handle.SetBPFFilter(filter)
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
					if dest_mac  == "" {
						saveHostMacAddr( packet )
					}
					pcapIncoming <- packet
				}
			}
        }(i)
    }
    go func() {
		    defer handle.Close()
			packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
			for {
				pcapPacket, _ := packetSource.NextPacket()
				pcapdQueue <- &pcapPacket
			}
	}()

    return pcapIncoming

}

func PollTimeoutRoutine( ipMeta * pState, timeoutQueue chan *packet_metadata, retransmitQueue chan *packet_metadata, 
	workers int, timeoutT int,  timeoutR int ) chan *packet_metadata  {

    TIMEOUT_T := time.Duration(timeoutT)*time.Second
    TIMEOUT_R := time.Duration(timeoutR)*time.Second

	timeoutIncoming := make(chan *packet_metadata, QUEUE_SIZE)
	//spawn off appropriate routines to poll from timeout & retransmit Queues at specified intervals
	timeoutAlg(  ipMeta, timeoutQueue, timeoutIncoming, TIMEOUT_T )
	timeoutAlg(  ipMeta, retransmitQueue, timeoutIncoming, TIMEOUT_R )

	return timeoutIncoming
}


//peek at front of q and sleep until processing
func timeoutAlg(  ipMeta * pState, queue chan *packet_metadata, timeoutIncoming chan *packet_metadata,
	timeout time.Duration) {

    go func() {
		tdif := time.Duration(timeout)
        for {
            select {
            case packet := <-queue:
				tdif = (time.Now()).Sub( packet.Timestamp )
				//if top of the Q is early, put routine to sleep until
				if tdif < timeout {
					//fmt.Println("sleeping:",timeout-tdif)
					time.Sleep(timeout-tdif)
				}

	            p, ok := ipMeta.find( packet )
                //if no longer in map
	            if !ok {
					//fmt.Println("not found")
                    continue
                }
                //if state hasnt changed
				if p.ExpectedRToLZR != packet.ExpectedRToLZR {
					//fmt.Println("state hasnt changed")
                    continue
                } else {
					//fmt.Println("will deal with")
                    timeoutIncoming <-packet
                }
            }
        }
    }()
}

// TimeoutQueueStuff TODO:need to move
func ConstructRetransmitQueue( workers int ) chan *packet_metadata {

    retransmitQueue := make(chan *packet_metadata, QUEUE_SIZE)
    return retransmitQueue
}



// TimeoutQueueStuff TODO:need to move
func ConstructTimeoutQueue( workers int ) chan *packet_metadata {

    timeoutQueue := make(chan *packet_metadata, QUEUE_SIZE)
    return timeoutQueue
}

