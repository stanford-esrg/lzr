package lzr

import (
	//"fmt"
	"log"
)


func closeConnection( packet *packet_metadata, ipMeta * pState, writingQueue chan *packet_metadata, write bool ) {

	//close connection
	rst := constructRST(packet)
	err := handle.WritePacketData(rst)
	if err != nil {
		log.Fatal(err)
	}
	//remove from state, we are done now
	packet = ipMeta.remove(packet)
	if write {
		writingQueue <- packet
	}
	return
}


func HandlePcap( opts *options, packet *packet_metadata, ipMeta * pState, timeoutQueue	chan *packet_metadata,
	retransmitQueue chan *packet_metadata, writingQueue chan *packet_metadata ) {


	//verify 
	if !(ipMeta.verifyScanningIP( packet )) {
		packet.incrementCounter()
		packet.updateTimestamp()
		packet.validationFail()
		timeoutQueue <-packet
		return
	}

	//for every ack received, mark as accepting data
	if (!packet.SYN) && packet.ACK {
		ipMeta.updateAck( packet )
	}

	 //exit condition
	 if len(packet.Data) > 0 {

		ipMeta.updateData( packet )

		// if not stopping here, send off to handle_expire
		if ForceAllHandshakes() {
			handleExpired( opts,packet, ipMeta, timeoutQueue, writingQueue )
			return
		}

		handshakeNum := ipMeta.getHandshake( packet )
		packet.syncHandshakeNum( handshakeNum )
		closeConnection( packet, ipMeta, writingQueue, true)
		return

	}
	//deal with closed connection 
	if packet.RST || packet.FIN {

		handleExpired( opts,packet, ipMeta, timeoutQueue, writingQueue )
		return

	 }
	 //for every ack received, mark as accepting data
	 if (!packet.SYN) && packet.ACK {
		 //add to map
		 packet.updateResponse(DATA)
		 packet.updateTimestamp()
		 ipMeta.update(packet)

		 //add to map
		 timeoutQueue <-packet
		 return
	}

	//for every s/a send the appropriate ack
	if packet.SYN && packet.ACK {

		if packet.HyperACKtive {
			closeConnection( packet, ipMeta, writingQueue, true )
			//remove the non hyperacktive state too
			packet.HyperACKtive = false
			closeConnection( packet, ipMeta, writingQueue, false)
			return
		}
		SendAck( opts, packet, ipMeta, timeoutQueue, retransmitQueue, writingQueue )

	}

}

