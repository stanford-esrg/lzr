package main

import (
    "fmt"
    "log"
)





func handleTimeout( packet packet_metadata, ipMeta * pState, timeoutQueue * chan packet_metadata, f *output_file ) {

	//verify that it wasnt already taken care of
	if !(ipMeta.verifyScanningIP( &packet )) {
	    return
	}

    //if it just acked, for now just write with empty data to file 
    //and do not requeue...will requeue later...

    //send again with just data no ack
    fmt.Println(packet)
    if packet.ExpectedR == ACK {
        if packet.Counter < 1 {
            packet.incrementCounter()
            data := constructData(packet,true,true)
            fmt.Println(data)
            err = handle.WritePacketData(data)
            if err != nil {
                log.Fatal(err)
            }
		    packet.updateTimestamp()
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

