/*
Copyright 2020 The Board of Trustees of The Leland Stanford Junior University

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package lzr

import (
    "log"
    //"fmt"
)



func SendAck( opts *options, synack  *packet_metadata, ipMeta * pState,
timeoutQueue  chan *packet_metadata, retransmitQueue chan *packet_metadata,
writingQueue  chan *packet_metadata, toACK bool, toPUSH bool, expectedResponse string ) {


	if synack.windowZero() {
		if !(RecordOnlyData()) {
			//not a real s/a
			writingQueue <- synack
		} else {
			addToSummary(synack)
		}
		return
	}

	//grab which handshake
	handshakeNum := ipMeta.getHandshake(synack)
	handshake, _ := GetHandshake( opts.Handshakes[ handshakeNum ] )

	//Send Ack with Data
	ack, payload := constructData( handshake, synack, toACK, toPUSH )//true, false )
	//add to map
	synack.updateResponse( expectedResponse )//ACK )
	synack.updateResponseL( payload )
	synack.updateTimestamp()
	ipMeta.update( synack )
	err := handle.WritePacketData(ack)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	synack.updateTimestamp()
	retransmitQueue <-synack
	return

}


