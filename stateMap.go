package main

type pState struct {

	IPmap	map[string]packet_metadata

}

//TODO: move the ipStateMap stuff to own file

/* keeps state by storing the packet that was sent 
 * and within the packet stores the expected response */
func constructPacketStateMap() pState {

	ipMeta := pState{
        IPmap: make( map[string]packet_metadata ),
    }
    return ipMeta
}


func (ipMeta * pState) metaContains(p * packet_metadata) bool {
	_, ok := (*ipMeta).IPmap[p.Saddr]
	if !ok {
		return false
	}
    return true
}

func (ipMeta * pState) update( packet packet_metadata ) {

	(*ipMeta).IPmap[packet.Saddr] = packet

}

func (ipMeta * pState) add( packet packet_metadata) {
        (*ipMeta).IPmap[packet.Saddr] = packet
}

func (ipMeta * pState) remove( packet packet_metadata) {
        delete( (*ipMeta).IPmap, packet.Saddr )
}

func ( ipMeta * pState ) verifyScanningIP( pRecv *packet_metadata ) bool {


	//first check that IP itself is being scanned
	pZMap, ok := (*ipMeta).IPmap[pRecv.Saddr]
	if !ok {
		return false
	}
	//second check that 4-tuple matches
	if (( pZMap.Saddr == pRecv.Saddr ) && (pZMap.Dport == pRecv.Dport) &&
		(pZMap.Sport == pRecv.Sport)) {
		return true
	}
	//TODO: check seq & ack and check state that we expect(?)

	return false

}


