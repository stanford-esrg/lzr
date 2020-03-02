package main

import (
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "fmt"
)





func handlePcap( packet gopacket.Packet, ipMeta * map[string]packet_metadata, timeoutQueue * chan packet_metadata ) {

        tcpLayer := packet.Layer(layers.LayerTypeTCP)
        if tcpLayer != nil {
            tcp, _ := tcpLayer.(*layers.TCP)
            ipLayer := packet.Layer(layers.LayerTypeIPv4)
            ip, _ := ipLayer.(*layers.IPv4)

            packet := *NewPacket(ip,tcp)
			//verify 
			if !(packet.verifyScanningIP( ipMeta )) {
				return
			}
			fmt.Println(packet)

            //for every ack received, mark as accepting data
            if (!tcp.SYN) && tcp.ACK {
                //TODO: do something with data
                fmt.Println(tcp.Payload)
                fmt.Println("acked")
                //close connection
                rst := constructRST(packet)
                err = handle.WritePacketData(rst)
            }
        }


}

