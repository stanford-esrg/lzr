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
    "log"
)


func HandleTimeout( opts *options, packet *packet_metadata, ipMeta * pState,
	timeoutQueue chan *packet_metadata, retransmitQueue  chan *packet_metadata,
	writingQueue  chan packet_metadata ) {

	//fmt.Println("Handling timeout")
	//fmt.Println(packet)
    //if packet has already been dealt with, return
    if !ipMeta.metaContains( packet ) {
		//fmt.Println("packet has already been dealt with in handle timeout! returning!")
        return
    }


    //send again with just data (not apart of handshake)
    if (packet.ExpectedRToLZR == DATA || packet.ExpectedRToLZR == ACK) {
        if packet.Counter < opts.RetransmitNum {
            packet.incrementCounter()


			//grab which handshake
			handshakeNum := ipMeta.getHandshake( packet )
			handshake, _ := GetHandshake( opts.Handshakes[ handshakeNum ] )

			//if packet counter is 0 then dont specify the push flag just yet
			dataPacket,payload := constructData( handshake, packet,true,!(packet.Counter  == 0))

            err = handle.WritePacketData( dataPacket )
            if err != nil {
                log.Fatal(err)
            }
		    packet.updateTimestamp()
			packet.updateResponseL( payload )
		    ipMeta.update( packet )

		    packet.updateTimestamp()
            timeoutQueue <- packet
	        return
        }
	}

	//this handshake timed-out 
	handleExpired( opts, packet, ipMeta, timeoutQueue, writingQueue )

    return

}

