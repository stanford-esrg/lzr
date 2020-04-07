package lzr

import (
    "log"
    //"fmt"
)



func SendAck( handshakes []string, synack  packet_metadata, ipMeta * pState, timeoutQueue * chan packet_metadata, 
    writingQueue * chan packet_metadata ) {

	//TODO: check that ip_metadata contains what we want (saddr,seq,ack,window)
	if synack.windowZero() {
		//not a real s/a
		*writingQueue <- synack
		return
	}

	//grab which handshake
	handshake := GetHandshake( handshakes[ synack.getHandshakeNum() ] )

	//Send Ack with Data
	ack, payload := constructData( handshake, synack, true, false )

	//add to map
	synack.updateResponse( ACK )
	synack.updateResponseL( payload )
	synack.updateTimestamp()
	ipMeta.update( &synack )
	err := handle.WritePacketData(ack)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	synack.updateTimestamp()
	*timeoutQueue <-synack
	return

}


