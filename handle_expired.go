package lzr

import (
	//"fmt"
)

func handleExpired( opts *options, packet * packet_metadata, ipMeta * pState, 
	timeoutQueue chan *packet_metadata, writingQueue chan *packet_metadata ) {

	// first close the existing connection unless
	// its already been terminated
	if !( packet.RST && !packet.ACK) {

		rst := constructRST( packet )
		err = handle.WritePacketData(rst)

	}

	//grab which handshake
	handshakeNum := ipMeta.getHandshake( packet )
	//if we are all out of handshakes to try, so sad. 
	if (packet.HyperACKtive  || (handshakeNum >= (len( opts.Handshakes ) - 1))){

		packet.syncHandshakeNum( handshakeNum )
		//remove from state, we are done now
		packet = ipMeta.remove(packet)
		if !packet.HyperACKtive {
			writingQueue <- packet
		}

	} else { // lets try another handshake

		//lets also filter for cananda-like things
		if ( handshakeNum == 0 &&  HyperACKtiveFiltering() ) {

			highPortPacket := createFilterPacket( packet )
			sendOffSyn( highPortPacket, ipMeta, timeoutQueue )

		}

		packet.updatePacketFlow()
        ipMeta.incHandshake( packet )
		sendOffSyn( packet, ipMeta, timeoutQueue )
	}

}

func sendOffSyn(packet * packet_metadata, ipMeta * pState,
    timeoutQueue chan *packet_metadata ) {

        packet.updateResponse( SYN_ACK )
        packet.updateTimestamp()
        ipMeta.update( packet )
        syn := constructSYN( packet )
        // send SYN packet if so and start the whole process again
        err = handle.WritePacketData(syn)
        if err != nil {
            panic(err)

        }
        //wait for a s/a
        packet.updateTimestamp()
        timeoutQueue <- packet

}
