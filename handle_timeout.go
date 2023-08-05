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


func HandleTimeout( opts *options, packet *packet_metadata, ipMeta * pState,
	timeoutQueue chan *packet_metadata, retransmitQueue  chan *packet_metadata,
	writingQueue  chan *packet_metadata ) {

	//fmt.Println("Handling timeout")
	//fmt.Println(packet)
    //if packet has already been dealt with, return
    if !ipMeta.metaContains( packet ) {
		//fmt.Println("packet has already been dealt with in handle timeout! returning!")
        return
    }


    //send again with just data (not apart of handshake)
    if ( packet.Counter < opts.RetransmitNum ) && !packet.HyperACKtive {
            packet.incrementCounter()

		if ( packet.ExpectedRToLZR == ACK || packet.ExpectedRToLZR == DATA ) {
			// pass in timeoutQ as retransmitQ to start the timeout clock
			SendAck( opts, packet, ipMeta, timeoutQueue, timeoutQueue, writingQueue,
					true, !(packet.Counter  == 0), packet.ExpectedRToLZR )
		}
		if ( packet.ExpectedRToLZR == SYN_ACK ) {
			SendSyn( packet, ipMeta, timeoutQueue )
		}
		return
	}

	//this handshake timed-out 
	handleExpired( opts, packet, ipMeta, timeoutQueue, writingQueue )

    return

}

