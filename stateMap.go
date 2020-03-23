package main

import (
    "sync"
    //"fmt"
)

type pState struct {

	IPmap	map[string]packet_metadata
    MLock    sync.RWMutex

}

/* keeps state by storing the packet that was received 
 * and within the packet stores the expected response 
 * storing received as to what was sent b/c want to know
 * perhaps need to wait some more 
 */
func constructPacketStateMap() pState {
    ipMeta := NewConcurrentMap()
    /*
	ipMeta := pState{
        IPmap: make( map[string]packet_metadata ),
        MLock: sync.RWMutex{},
    }*/
    return ipMeta
}


func (ipMeta * pState) metaContains( p * packet_metadata ) bool {
    ipMeta.MLock.RLock()
	_, ok := ipMeta.IPmap[p.Saddr]
    ipMeta.MLock.RUnlock()
	if !ok {
		return false
	}
    return true
}


//is Processing for goPackets
func (ipMeta * pState) isProcessing( p * packet_metadata ) ( bool,bool ) {
    ipMeta.MLock.RLock()
    defer ipMeta.MLock.RUnlock()
    p_out, ok := ipMeta.IPmap[p.Saddr]
    if !ok {
        return false,false
    }
    return true, p_out.Processing

}

func (ipMeta * pState) startProcessing( p * packet_metadata ) bool {

    ipMeta.MLock.Lock()
    defer ipMeta.MLock.Unlock()
    p_out, ok := ipMeta.IPmap[p.Saddr]
    if !ok {
        return false
    }
    p_out.startProcessing()
    return true

}

func (ipMeta * pState) finishProcessing( p * packet_metadata ) bool {

    ipMeta.MLock.RLock()
    defer ipMeta.MLock.RUnlock()
    p_out, ok := ipMeta.IPmap[p.Saddr]
    if !ok {
        return false
    }
    p_out.finishedProcessing()
    return true

}



func (ipMeta * pState) find(p * packet_metadata) ( *packet_metadata, bool ) {
    ipMeta.MLock.RLock()
	p_out, ok := ipMeta.IPmap[p.Saddr]
    ipMeta.MLock.RUnlock()
    return &p_out, ok
}
func (ipMeta * pState) update( packet * packet_metadata ) {
    ipMeta.MLock.Lock()
	ipMeta.IPmap[packet.Saddr] = *packet
    ipMeta.MLock.Unlock()
}

func (ipMeta * pState) incrementCounter( packet * packet_metadata ) bool {

    p_out, ok := ipMeta.find( packet )
    if !ok {
        return false
    }
    p_out.incrementCounter()
    ipMeta.update( p_out )
    return true
}


func (ipMeta * pState) remove( packet packet_metadata) {
    ipMeta.MLock.Lock()
    delete( ipMeta.IPmap, packet.Saddr )
    ipMeta.MLock.Unlock()
}

func ( ipMeta * pState ) verifyScanningIP( pRecv *packet_metadata ) bool {


	//first check that IP itself is being scanned
    ipMeta.MLock.RLock()
	pMap, ok := ipMeta.IPmap[pRecv.Saddr]
    ipMeta.MLock.RUnlock()
	if !ok {
		return false
	}
	//second check that 4-tuple matches
	//TODO: check seq & ack and check state that we expect(?)
	if (( pMap.Saddr == pRecv.Saddr ) && (pMap.Dport == pRecv.Dport) &&
    (pMap.Sport == pRecv.Sport) ) { // && (pRecv.Acknum == pMap.Seqnum + 1)) {

            /*
             fmt.Println("recv seq num:", pRecv.Seqnum)
             fmt.Println("stored seqnum: ", pMap.Seqnum)
             fmt.Println("recv ack num:", pRecv.Acknum)
             fmt.Println("stored acknum: ", pMap.Acknum)
             fmt.Println("received response length: ",len(pRecv.Data))
             fmt.Println("stored response length: ",pMap.LZRResponseL) 
            */
            /* if ( pRecv.Acknum == ( pMap.Acknum + pMap.LZRResponseL ) ) {
                 if ((pRecv.Seqnum == ( pMap.Seqnum )) || (pRecv.Seqnum == ( pMap.Seqnum + 1 ))) {
                    return true
                 } */
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
                 }
             }*/
             return true

    }

	return false

}


