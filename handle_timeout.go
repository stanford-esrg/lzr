package lzr

import (
    //"fmt"
    "log"
)


func HandleTimeout( handshake *Handshake, packet packet_metadata, ipMeta * pState, timeoutQueue * chan packet_metadata,
    writingQueue * chan packet_metadata, f *output_file ) {

    //if packet has already been dealt with, return
    if !ipMeta.metaContains( &packet ) {
        return
    }

    //send again with just data (not apart of handshake)
    if (packet.ExpectedRToLZR == DATA || packet.ExpectedRToLZR == ACK) {
        if packet.Counter < 2 {
            packet.incrementCounter()
			var dataPacket []byte
			if packet.Counter  == 0 {
				dataPacket,_ = constructData( handshake, packet,true,false )
			} else {
				dataPacket,_ = constructData( handshake, packet,true,true )
			}
            err = handle.WritePacketData( dataPacket )
            if err != nil {
                log.Fatal(err)
            }
		    packet.updateTimestamp()
		    ipMeta.update( &packet )

		    packet.updateTimestamp()
            *timeoutQueue <- packet
	        return
        }
	}
    //if pcap still has something relevant to this ip, put it back on timeout queue 
    //TODO: update timestamp or just to the back of Q is good enough??
    /*if packet.PCapTracker != 0 {
        *timeoutQueue <- packet
    }*/

    //else, we give up, just record. 
    //remove from state, we are done now
    ipMeta.remove(packet)
	//fmt.Println("timeout Removed <- ",string(packet.Saddr))
    *writingQueue <- packet
    //close connection
    rst := constructRST(packet)
    err = handle.WritePacketData(rst)
    return

}

