package main

import (
    //"sync"
    //"fmt"
)

/* keeps state by storing the packet that was received 
 * and within the packet stores the expected response 
 * storing received as to what was sent b/c want to know
 * perhaps need to wait some more 
 */
func constructPacketStateMap() pState {
    ipMeta := NewpState()
    return ipMeta
}


func (ipMeta * pState) metaContains( p * packet_metadata ) bool {
    return ipMeta.Has(p.Saddr)
}


func (ipMeta * pState) find(p * packet_metadata) ( *packet_metadata, bool ) {
    return ipMeta.Get(p.Saddr)
}

func (ipMeta * pState) update( packet * packet_metadata ) {
    ipMeta.Insert( packet.Saddr, packet )
}

func (ipMeta * pState) incrementCounter( packet * packet_metadata ) bool {

    p_out, ok := ipMeta.Get(packet.Saddr)
    if !ok {
        return false
    }
    p_out.incrementCounter()
    ipMeta.Insert( p_out.Saddr, p_out )
    return true
}


func (ipMeta * pState) remove( packet packet_metadata) {
    ipMeta.Remove(packet.Saddr)
    return
}

func ( ipMeta * pState ) verifyScanningIP( pRecv *packet_metadata ) bool {

	//first check that IP itself is being scanned
    pMap, ok := ipMeta.Get(pRecv.Saddr)
	if !ok {
		return false
	}
	//second check that 4-tuple matches
	//TODO: check seq & ack and check state that we expect(?)
	if (( pMap.Saddr == pRecv.Saddr ) && (pMap.Dport == pRecv.Dport) &&
    (pMap.Sport == pRecv.Sport) ) { // && (pRecv.Acknum == pMap.Seqnum + 1)) {

            if ( pRecv.Acknum == ( pMap.Acknum + pMap.LZRResponseL ) ) {
                 if ((pRecv.Seqnum == ( pMap.Seqnum )) || (pRecv.Seqnum == ( pMap.Seqnum + 1 ))) {
                    return true
                 }
/*                 //fmt.Println("ack passed")
                 if (len(pRecv.Data) > 0 ) {
                    if pRecv.Seqnum == ( pMap.Seqnum + 1) {
                        //fmt.Println("here")
                        return true
                    }
                 } else {
                    if pRecv.Seqnum == ( pMap.Seqnum  ){
                        return true
                    }
                 }*/
             }
             //return true

    }

             /*fmt.Println(pMap.Saddr, "====")
             fmt.Println("recv seq num:", pRecv.Seqnum)
             fmt.Println("stored seqnum: ", pMap.Seqnum)
             fmt.Println("recv ack num:", pRecv.Acknum)
             fmt.Println("stored acknum: ", pMap.Acknum)
             fmt.Println("received response length: ",len(pRecv.Data))
             fmt.Println("stored response length: ",pMap.LZRResponseL) 
             fmt.Println(pMap.Saddr ,"====")*/
	return false

}


