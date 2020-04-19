package lzr

import (
    //"fmt"
    "log"
)





func HandlePcap( handshakes []string, packet *packet_metadata, ipMeta * pState, timeoutQueue * chan *packet_metadata, 
    writingQueue * chan *packet_metadata ) {


    //packet.PCapTracker -= 1

    //verify 
	if !(ipMeta.verifyScanningIP( packet )) {
        packet.incrementCounter()
		packet.updateTimestamp()
        packet.validationFail()
        *timeoutQueue <-packet
		return
	}

     //exit condition
     if len(packet.Data) > 0 {
		handshakeNum := ipMeta.getHandshake( packet )
		packet.syncHandshakeNum( handshakeNum )
        packet.fingerprintData()
        *writingQueue <- packet
        //close connection
        rst := constructRST(packet)
        err := handle.WritePacketData(rst)
        if err != nil {
            log.Fatal(err)
        }
        //remove from state, we are done now
        ipMeta.remove(packet)
        return

    }
    //deal with closed connection 
    if packet.RST || packet.FIN {

		handleExpired( handshakes,packet, ipMeta, writingQueue )
		return

     }
     //for every ack received, mark as accepting data
     if (!packet.SYN) && packet.ACK {

		 //add to map
		 packet.updateResponse(DATA)
		 packet.updateTimestamp()
		 ipMeta.update(packet)

		 //add to map
         *timeoutQueue <-packet
		  return
    }

	//for every s/a send the appropriate ack
	if packet.SYN && packet.ACK {

		SendAck( handshakes, packet, ipMeta, timeoutQueue, writingQueue )

	}

}

