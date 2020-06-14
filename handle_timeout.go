package lzr

import (
    //"fmt"
    "log"
)


func HandleTimeout( opts *options, packet *packet_metadata, ipMeta * pState,
	timeoutQueue chan *packet_metadata, retransmitQueue  chan *packet_metadata,
	writingQueue  chan packet_metadata ) {

    //if packet has already been dealt with, return
    if !ipMeta.metaContains( packet ) {
        return
    }


    //send again with just data (not apart of handshake)
    if (packet.ExpectedRToLZR == DATA || packet.ExpectedRToLZR == ACK) {
        if packet.Counter < opts.RetransmitNum {
            packet.incrementCounter()


			//grab which handshake
			handshakeNum := ipMeta.getHandshake( packet )
			handshake, _ := GetHandshake( opts.Handshakes[ handshakeNum ] )

			//if packet counter is 0 then dont specify the push flag just yet
			dataPacket,payload := constructData( handshake, packet,true,!(packet.Counter  == 0))

            err = handle.WritePacketData( dataPacket )
            if err != nil {
                log.Fatal(err)
            }
		    packet.updateTimestamp()
			packet.updateResponseL( payload )
		    ipMeta.update( packet )

		    packet.updateTimestamp()
            timeoutQueue <- packet
	        return
        }
	}

	//this handshake timed-out 
	handleExpired( opts, packet, ipMeta, timeoutQueue, writingQueue )

    return

}

