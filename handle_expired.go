package lzr

import (
	//"fmt"
)

func handleExpired( opts *options, packet * packet_metadata, ipMeta * pState,
	timeoutQueue chan *packet_metadata, writingQueue chan packet_metadata ) {


	// first close the existing connection unless
	// its already been terminated
	if !( packet.RST && !packet.ACK ) {

		rst := constructRST( packet )
		_ = handle.WritePacketData( rst )

	}

	//grab which handshake
	handshakeNum := ipMeta.getHandshake( packet )

	//if we are all not trying anymore handshakes, so sad. 
	//if ( packet.ExpectedRToLZR == SYN_ACK ||
	if ( packet.HyperACKtive  || (handshakeNum >= (len( opts.Handshakes ) - 1)) ){

		packet.syncHandshakeNum( handshakeNum )

		//document failure if its a handshake response that hasnt succeeded before
		if !packet.HyperACKtive && !( ForceAllHandshakes() && ipMeta.getData( packet ) && len(packet.Data) == 0 ) {
			writingQueue <- *packet
		}

		//remove from state, we are done now
		packet = ipMeta.remove( packet )
		if HyperACKtiveFiltering() {
			packet.HyperACKtive = true
			packet = ipMeta.remove( packet )
		}
	} else { // lets try another handshake


		//record all succesful fingerprints if forcing all handshakes
		if ForceAllHandshakes() && len(packet.Data) > 0 {
			packet.syncHandshakeNum( handshakeNum )
			writingQueue <- *packet
		}

		packet.updatePacketFlow()
		ipMeta.incHandshake( packet )
		sendOffSyn( packet, ipMeta, timeoutQueue )

		//lets also filter for cananda-like things
		if ( handshakeNum == 0 &&  HyperACKtiveFiltering() ) {
			for i := 0; i < getNumFilters(); i++ {
				highPortPacket := createFilterPacket( packet )
				sendOffSyn( highPortPacket, ipMeta, timeoutQueue )
				ipMeta.incHandshake( highPortPacket )
				ipMeta.setHyperACKtiveStatus( highPortPacket )

				ipMeta.setParentSport( highPortPacket, packet.Sport )
				ipMeta.FinishProcessing( highPortPacket )
			}
		}

	}

}

func sendOffSyn(packet * packet_metadata, ipMeta * pState,
	timeoutQueue chan *packet_metadata ) {

		packet.updateResponse( SYN_ACK )
		packet.updateTimestamp()
		ipMeta.update( packet )
		syn := constructSYN( packet )
		// send SYN packet if so and start the whole process again
		err := handle.WritePacketData(syn)
		if err != nil {
			panic(err)

		}
		//wait for a s/a
		packet.updateTimestamp()
		ipMeta.FinishProcessing( packet )
		timeoutQueue <- packet

}
