package main

import (
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    //"fmt"
)





func handlePcap( packet gopacket.Packet, ipMeta * pState, timeoutQueue * chan packet_metadata, f *output_file ) {

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
            //TODO : change to use packet and not tcp
            //for every ack received, mark as accepting data
            if (!tcp.SYN) && tcp.ACK {

                //exit condition
                if len(tcp.Payload) > 0 {

                    //remove from state, we are done now
                    ipMeta.remove(packet)
                    //TODO: do something with data
                    //like fingerprint
                    packet.updateData(string(tcp.Payload))
                    f.record(packet)
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
            //for every closed connection, record
            if tcp.RST || tcp.FIN {

                

            }

        }


}

