package main

import (
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    //"fmt"
    "log"
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
            //for every closed connection, record
            if packet.RST || packet.FIN {

                // do not remove packet due to race condition
                //ipMeta.remove(packet)
                f.record(packet)
                //close connection
                if packet.FIN {
                    rst := constructRST(packet)
                    err = handle.WritePacketData(rst)
                }
                return

            }
            //for every ack received, mark as accepting data
            if (!packet.SYN) && packet.ACK {

                //exit condition
                if len(packet.Data) > 0 {

                    //remove from state, we are done now
                    ipMeta.remove(packet)
                    //TODO: do something with data
                    //like fingerprint
                    f.record(packet)
                    //close connection
                    rst := constructRST(packet)
                    err = handle.WritePacketData(rst)
                    if err != nil {
                       log.Fatal(err)
                    }
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

