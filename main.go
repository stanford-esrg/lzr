package main

import (
    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
    "log"
    "io"
    "bufio"
    "os"
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



func main() {

	//initalize
	ipMeta := make( map[string]packet_metadata )

    //read in config 
    //port := parse()

    zmapIncoming := constructZMapRoutine()
    pcapIncoming := constructPcapRoutine()

	//read from both zmap and pcap
	for {
		select {
			case input := <-zmapIncoming:
				ackZMap( input, &ipMeta )
			case input := <-pcapIncoming:
				handlePcap( input, &ipMeta )
			default:
				//continue to non-blocking poll
		}
	}

} //end of main
