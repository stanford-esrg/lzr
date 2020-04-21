package lzr

import (
	//"fmt"
)

func handleExpired( handshakes []string, packet * packet_metadata, ipMeta * pState, writingQueue chan *packet_metadata ) {

	// first close the existing connection unless
	// its already been terminated
	if !packet.RST {

		rst := constructRST( packet )
		err = handle.WritePacketData(rst)

	}

	//grab which handshake
	handshakeNum := ipMeta.getHandshake( packet )

	//if we are all out of handshakes to try, so sad. 
	if handshakeNum >= (len( handshakes ) - 1){
		packet.syncHandshakeNum( handshakeNum )
		//remove from state, we are done now
		packet = ipMeta.remove(packet)
		writingQueue <- packet

	} else { // lets try another handshake

		packet.updatePacketFlow()
		ipMeta.incHandshake( packet )
		ipMeta.update( packet )
		syn := constructSYN( packet )
		// send SYN packet if so and start the whole process again
		err = handle.WritePacketData(syn)
		if err != nil {
			panic(err)

		}
	}

}
