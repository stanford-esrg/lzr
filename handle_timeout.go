package main

import (
    "fmt"
    "log"
)





func handleTimeout( packet packet_metadata, ipMeta * pState, timeoutQueue * chan packet_metadata, f *output_file ) {

	//verify that it wasnt already taken care of

    _,correctnum := ipMeta.verifyScanningIP( &packet )
    //TODO make esception for correctnum
	if !(correctnum) {
	    return
	}

    if len(packet.Data) > 0 {
        fmt.Println("DATA APPEARED")
        f.record(packet)
        //close connection
        rst := constructRST(packet)
        err = handle.WritePacketData(rst)
    }
    //if it just acked, for now just write with empty data to file 
    //and do not requeue...will requeue later...

    //send again with just data no ack
    if packet.ExpectedR == ACK {
        if packet.Counter < 1 {
            packet.incrementCounter()
            data := getData( string(packet.Saddr) )
            dataPacket := constructData( packet,data, true,true )
            err = handle.WritePacketData( dataPacket )
            if err != nil {
                log.Fatal(err)
            }
		    packet.updateTimestamp()
		    ipMeta.update(packet)
            *timeoutQueue <- packet
            return
        }
    }

    //remove from state, we are done now
    ipMeta.remove(packet)
    f.record(packet)
    //close connection
    rst := constructRST(packet)
    err = handle.WritePacketData(rst)
    return

}

