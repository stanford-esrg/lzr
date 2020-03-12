package main

import (
    "sync"
    "fmt"
)

type pState struct {

	IPmap	map[string]packet_metadata
    MLock    sync.RWMutex

}

//TODO: move the ipStateMap stuff to own file

/* keeps state by storing the packet that was received 
 * and within the packet stores the expected response 
 * storing received as to what was sent b/c want to know
 * perhaps need to wait some more 
 */
func constructPacketStateMap() pState {

	ipMeta := pState{
        IPmap: make( map[string]packet_metadata ),
        MLock: sync.RWMutex{},
    }
    return ipMeta
}


func (ipMeta * pState) metaContains(p * packet_metadata) bool {
    ipMeta.MLock.RLock()
	_, ok := ipMeta.IPmap[p.Saddr]
    ipMeta.MLock.RUnlock()
	if !ok {
		return false
	}
    return true
}

func (ipMeta * pState) find(p * packet_metadata) ( *packet_metadata, bool ) {
    ipMeta.MLock.RLock()
	p_out, ok := ipMeta.IPmap[p.Saddr]
    ipMeta.MLock.RUnlock()
    return &p_out, ok
}
func (ipMeta * pState) update( packet packet_metadata ) {
    ipMeta.MLock.Lock()
	ipMeta.IPmap[packet.Saddr] = packet
    ipMeta.MLock.Unlock()
}


func (ipMeta * pState) remove( packet packet_metadata) {
    ipMeta.MLock.Lock()
    delete( ipMeta.IPmap, packet.Saddr )
    ipMeta.MLock.Unlock()
}

//returns if its an ip totally not worth considering or potentially just out of order
func ( ipMeta * pState ) verifyScanningIP( pRecv *packet_metadata ) (bool,bool) {


	//first check that IP itself is being scanned
    ipMeta.MLock.RLock()
	pMap, ok := ipMeta.IPmap[pRecv.Saddr]
    ipMeta.MLock.RUnlock()
	if !ok {
		return false,false
	}
	//second check that 4-tuple matches
	if ( ( pMap.Saddr == pRecv.Saddr ) && ( pMap.Dport == pRecv.Dport ) &&
    ( pMap.Sport == pRecv.Sport ) ) {
        if ( pMap.ResponseL > 0 ) {
            fmt.Println("recv ack num:", pRecv.Acknum)
            fmt.Println("stored seqnum: ", pMap.Seqnum)
            fmt.Println("stored response length: ",pMap.ResponseL)
            if ( pRecv.Acknum == ( pMap.Acknum + pMap.ResponseL ) ) {
                fmt.Println("TRUE_ RESPONSE")
		        return true,true
            }
        } else {
            if ( pRecv.Acknum == ( pMap.Acknum ) ) {
                fmt.Println("TRUE_SEQ+1")
                return true,true
            }
        }
	}
    //seq and ack dont _totally_ match up
	return false,true

}


