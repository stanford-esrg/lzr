package main

import (
    "log"
    "fmt"
)



func ackZMap( synack  packet_metadata, ipMeta * pState, timeoutQueue * chan packet_metadata, 
    writingQueue * chan packet_metadata, f *output_file ) {

        //TODO: check that ip_metadata contains what we want (saddr,seq,ack,window)

        fmt.Println(synack.Saddr)
        if synack.windowZero() {
            //not a real s/a
            *writingQueue <- synack
            return
        }

        //Send Ack with Data
        data := getData( string(synack.Saddr) )
        ack := constructData( synack, data, true, false )

		//add to map
		synack.updateResponse( ACK )
        synack.updateResponseL( data )
		synack.updateTimestamp()
		ipMeta.update( synack )
        err = handle.WritePacketData(ack)
        if err != nil {
            //panic(err)
            log.Fatal(err)
        }
        *timeoutQueue <-synack
		return

}


