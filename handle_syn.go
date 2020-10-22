package lzr

import (
	//"fmt"
)

func SendOffSyn(packet * packet_metadata, ipMeta * pState,
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
