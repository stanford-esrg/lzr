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
	//"fmt"
)

func SendSyn(packet * packet_metadata, ipMeta * pState,
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
