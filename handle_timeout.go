package main

import (
    //"fmt"
    "log"
)


func handleTimeout( packet packet_metadata, ipMeta * pState, timeoutQueue * chan packet_metadata,
    writingQueue * chan packet_metadata, f *output_file ) {

    //if packet has already been dealt with, return
    if !ipMeta.metaContains( &packet ) {
        return
    }

   //TODO: not sure this is what i want
   //this is if seq/ack dont match up b/c server refuses to 
   //ack our data. 
   if packet.getValidationFail() {
       //if we waited a bit and validation seems to still fail, just record
       // and move on
       //2 times timout is arbitrary
       if packet.Counter >= 2 {

           //remove from state, we are done now
           ipMeta.remove(packet)
           *writingQueue <- packet
           //close connection
           rst := constructRST(packet)
           err = handle.WritePacketData(rst)
       }
       //take note that validation failed again
       packet.incrementCounter()
   }

	//verify that it wasnt already taken care of
	if !(ipMeta.verifyScanningIP( &packet )) {


        packet.validationFail()
        //send again with just data (not apart of handshake)
        if packet.ExpectedRToLZR == ACK {
            if packet.Counter < 1 {
                packet.incrementCounter()
                data := getData( string(packet.Saddr) )
                dataPacket := constructData( packet,data, true,true )
                err = handle.WritePacketData( dataPacket )
                if err != nil {
                    log.Fatal(err)
                }
		        packet.updateTimestamp()
		        ipMeta.update( &packet )
            }
       }

		packet.updateTimestamp()
        *timeoutQueue <- packet
	    return
	}

    //if it just acked, for now just write with empty data to file 
    //and do not requeue...will requeue later...


    //remove from state, we are done now
    ipMeta.remove(packet)
    *writingQueue <- packet
    //close connection
    rst := constructRST(packet)
    err = handle.WritePacketData(rst)
    return

}

