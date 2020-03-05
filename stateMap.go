package main

import (
    "sync"
)

type pState struct {

	IPmap	map[string]packet_metadata
    MLock    sync.RWMutex

}

//TODO: move the ipStateMap stuff to own file

/* keeps state by storing the packet that was sent 
 * and within the packet stores the expected response */
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

func ( ipMeta * pState ) verifyScanningIP( pRecv *packet_metadata ) bool {


	//first check that IP itself is being scanned
    ipMeta.MLock.RLock()
	pMap, ok := ipMeta.IPmap[pRecv.Saddr]
    ipMeta.MLock.RUnlock()
	if !ok {
		return false
	}
	//second check that 4-tuple matches
	if (( pMap.Saddr == pRecv.Saddr ) && (pMap.Dport == pRecv.Dport) &&
		(pMap.Sport == pRecv.Sport)) {
		return true
	}
	//TODO: check seq & ack and check state that we expect(?)

	return false

}


