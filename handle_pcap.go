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


func closeConnection( packet *packet_metadata, ipMeta * pState, writingQueue chan *packet_metadata, write bool, ackingFirewall bool ) {

	//close connection
	rst := constructRST(packet)
	err := handle.WritePacketData(rst)
	if err != nil {
		log.Fatal(err)
	}
	//remove from state, we are done now
	packet = ipMeta.remove(packet)
	if write {
		packet.setHyperACKtive(ackingFirewall)

        if !(!(packet.hasData()) && RecordOnlyData()){
             writingQueue <- packet
        } else {
            addToSummary(packet)
        }
	}
	return
}


func HandlePcap( opts *options, packet *packet_metadata, ipMeta * pState, timeoutQueue	chan *packet_metadata,
	retransmitQueue chan *packet_metadata, writingQueue chan *packet_metadata ) {


	//verify
	verified := ipMeta.verifyScanningIP( packet )
	if !verified {
		packet.incrementCounter()
		packet.updateTimestamp()
		packet.validationFail()
		timeoutQueue <-packet
		return
	}

	isHyperACKtive := ipMeta.getHyperACKtiveStatus( packet )
	handshakeNum := ipMeta.getHandshake( packet )

	//for every ack received, mark as accepting data
	if (!packet.SYN) && packet.ACK {
		ipMeta.updateAck( packet )
	}
	 //exit condition
	 if len(packet.Data) > 0 {
		packet.updateResponse(DATA)
		ipMeta.updateData( packet )

		// if not stopping here, send off to handle_expire
		if ForceAllHandshakes() {
			handleExpired( opts,packet, ipMeta, timeoutQueue, writingQueue )
			return
		}

		packet.syncHandshakeNum( handshakeNum )

		closeConnection( packet, ipMeta, writingQueue, true,  isHyperACKtive)
		return

	}
	//deal with closed connection 
	if packet.RST || packet.FIN {

		handleExpired( opts,packet, ipMeta, timeoutQueue, writingQueue )
		return

	 }


	//checking if max filter syn acks reached
	//( filterACKs + original ACK + this ack)
     if handshakeNum == 1 && HyperACKtiveFiltering() && !isHyperACKtive {
			//fmt.Println( ipMeta.getEphemeralRespNum( packet ) )
			//fmt.Println(getNumFilters())
            if ipMeta.getEphemeralRespNum( packet )   > getNumFilters() {
                closeConnection( packet, ipMeta, writingQueue, true, true)
				return
            }
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

		if  handshakeNum == 1 && HyperACKtiveFiltering() {

			//just close and record
			if isHyperACKtive {

                parentSport := ipMeta.getParentSport( packet )

				ipMeta.incEphemeralResp( packet, parentSport )
				closeConnection( packet, ipMeta, writingQueue, false, isHyperACKtive)
				return
			} else {
				ipMeta.incEphemeralResp( packet, packet.Sport )
			}
		}
		toACK := true
		toPUSH := false
		SendAck( opts, packet, ipMeta, timeoutQueue, retransmitQueue, writingQueue,
				toACK, toPUSH, ACK )
		return
	}

}

