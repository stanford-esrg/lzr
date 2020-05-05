package lzr

import (
	"fmt"
)

/* keeps state by storing the packet that was received 
 * and within the packet stores the expected response 
 * storing received as to what was sent b/c want to know
 * perhaps need to wait some more 
 */
func ConstructPacketStateMap( opts *options ) pState {
    ipMeta := NewpState()
    return ipMeta
}


func constructKey( packet *packet_metadata ) string {

	return packet.Saddr+packet.Daddr+string(packet.Sport)+string(packet.Dport)

}

func (ipMeta * pState) metaContains( p * packet_metadata ) bool {
    return ipMeta.Has(p.Key)
}


func (ipMeta * pState) find(p * packet_metadata) ( *packet_metadata, bool ) {
    ps, ok := ipMeta.Get(p.Key)
	if ok {
		return ps.Packet, ok
	}
    return nil,ok
}

func (ipMeta * pState) update( p * packet_metadata ) {

    ps, ok := ipMeta.Get(p.Key)
	if !ok {
		ps = &packet_state {
			Packet: p,
			Ack: false,
			HandshakeNum: 0,
		}
	} else {
		ps.Packet = p
	}
    ipMeta.Insert( p.Key, ps )
}


func (ipMeta * pState) incHandshake( p * packet_metadata ) bool {
    ps, ok := ipMeta.Get(p.Key)
	if ok {
		ps.HandshakeNum += 1
		ipMeta.Insert( p.Key, ps )
	}
	return ok
}

func (ipMeta * pState) updateAck( p * packet_metadata ) bool {
    ps, ok := ipMeta.Get(p.Key)
    if ok {
        ps.Ack = true
        ipMeta.Insert( p.Key, ps )
    }
    return ok
}

func (ipMeta * pState) getAck( p * packet_metadata ) bool {
    ps, ok := ipMeta.Get(p.Key)
    if ok {
        return ps.Ack
    }
    return false
}

func (ipMeta * pState) getHandshake( p * packet_metadata ) int {
    ps, ok := ipMeta.Get(p.Key)
    if ok {
        return ps.HandshakeNum
    }
    return 0
}

func (ipMeta * pState) incrementCounter( p * packet_metadata ) bool {

    ps, ok := ipMeta.Get(p.Key)
    if !ok {
        return false
    }
	ps.Packet.incrementCounter()
    ipMeta.Insert( ps.Packet.Key, ps )
    return true

}


func (ipMeta * pState) remove( packet *packet_metadata ) *packet_metadata {
	packet.ACKed = ipMeta.getAck( packet )
    ipMeta.Remove( packet.Key )
    return packet
}

func ( ipMeta * pState ) verifyScanningIP( pRecv *packet_metadata ) bool {

	//first check that IP itself is being scanned
    ps, ok := ipMeta.Get(pRecv.Key)
	if !ok {
		return false
	}
	pMap := ps.Packet

	//second check that 4-tuple matches
	//TODO: check seq & ack and check state that we expect(?)
	if (( pMap.Saddr == pRecv.Saddr ) && (pMap.Dport == pRecv.Dport) &&
    (pMap.Sport == pRecv.Sport) ) { // && (pRecv.Acknum == pMap.Seqnum + 1)) {

		if pRecv.SYN && pRecv.ACK {
			if ( pRecv.Acknum == pMap.Seqnum + 1 ) {
				return true
			}
		} else {

			if ((pRecv.Seqnum == ( pMap.Seqnum )) || (pRecv.Seqnum == ( pMap.Seqnum + 1 ))) {

				if ( pRecv.Acknum == ( pMap.Acknum + pMap.LZRResponseL ) ) {
					return true
				}
			}
		}
	}
	if DebugOn() {
		fmt.Println(pMap.Saddr, "====")
		fmt.Println("recv seq num:", pRecv.Seqnum)
		fmt.Println("stored seqnum: ", pMap.Seqnum)
		fmt.Println("recv ack num:", pRecv.Acknum)
		fmt.Println("stored acknum: ", pMap.Acknum)
		fmt.Println("received response length: ",len(pRecv.Data))
		fmt.Println("stored response length: ",pMap.LZRResponseL)
		fmt.Println(pMap.Saddr ,"====")
	}
	return false

}


