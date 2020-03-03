package main

import (
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "fmt"
)





func handlePcap( packet gopacket.Packet, ipMeta * pState, timeoutQueue * chan packet_metadata ) {

        tcpLayer := packet.Layer(layers.LayerTypeTCP)
        if tcpLayer != nil {
            tcp, _ := tcpLayer.(*layers.TCP)
            ipLayer := packet.Layer(layers.LayerTypeIPv4)
            ip, _ := ipLayer.(*layers.IPv4)

            packet := *NewPacket(ip,tcp)
			//verify 
			if !(ipMeta.verifyScanningIP( &packet )) {
				return
			}
			fmt.Println(packet)

            //for every ack received, mark as accepting data
            if (!tcp.SYN) && tcp.ACK {

                //exit condition
                if len(tcp.Payload) > 0 {

                    //remove from state, we are done now
                    ipMeta.remove(packet)
                    //TODO: do something with data
                    fmt.Println(tcp.Payload)
                    //close connection
                    rst := constructRST(packet)
                    err = handle.WritePacketData(rst)
                    return

                }

		        //add to map
		        packet.updateState(DATA)
		        packet.updateTimestamp()
		        ipMeta.update(packet)

		        //add to map
                *timeoutQueue <-packet

            }
        }


}

