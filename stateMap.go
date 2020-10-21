package lzr

import (
	"strconv"
	"fmt"
)

/* keeps state by storing the packet that was received 
 * and within the packet stores the expected response. 
 * storing received as to what was sent b/c want to know
 * perhaps need to wait some more 
 */
func ConstructPacketStateMap( opts *options ) pState {
	ipMeta := NewpState()
	return ipMeta
}


func constructKey( packet *packet_metadata ) string {
	return packet.Saddr + ":" + strconv.Itoa(packet.Sport)
}

func constructParentKey( packet *packet_metadata, parentSport int ) string {
    return packet.Saddr + ":" + strconv.Itoa(parentSport)
}


func (ipMeta * pState) metaContains( p * packet_metadata ) bool {

	pKey := constructKey(p)
	return ipMeta.Has(pKey)

}


func (ipMeta * pState) find(p * packet_metadata) ( *packet_metadata, bool ) {
	pKey := constructKey(p)
	ps, ok := ipMeta.Get(pKey)
	if ok {
		return ps.Packet, ok
	}
	return nil,ok
}

func (ipMeta * pState) update( p * packet_metadata ) {

	pKey := constructKey(p)
	ps, ok := ipMeta.Get(pKey)
	if !ok {
		ps = &packet_state {
			Packet: p,
			Ack: false,
			HandshakeNum: 0,
		}
	} else {
		ps.Packet = p
	}
	ipMeta.Insert( pKey, ps )
}


func (ipMeta * pState) incHandshake( p * packet_metadata ) bool {
	pKey := constructKey(p)
	ps, ok := ipMeta.Get(pKey)
	if ok {
		ps.HandshakeNum += 1
		ipMeta.Insert( pKey, ps )
	}
	return ok
}

func (ipMeta * pState) updateAck( p * packet_metadata ) bool {
	pKey := constructKey(p)
	ps, ok := ipMeta.Get(pKey)
	if ok {
		ps.Ack = true
		ipMeta.Insert( pKey, ps )
	}
	return ok
}

func (ipMeta * pState) getAck( p * packet_metadata ) bool {
	pKey := constructKey(p)
	ps, ok := ipMeta.Get(pKey)
	if ok {
		return ps.Ack
	}
	return false
}

func (ipMeta * pState) incEphemeralResp( p * packet_metadata, sport int) bool {
    pKey := constructParentKey(p, sport)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        ps.EphemeralRespNum += 1
        ipMeta.Insert( pKey, ps )
    }
    return ok
}


func (ipMeta * pState) getEphemeralRespNum( p * packet_metadata ) int {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        return ps.EphemeralRespNum
    }
    return 0
}

func (ipMeta * pState) getHyperACKtiveStatus( p * packet_metadata ) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        return ps.HyperACKtive
    }
    return false
}

func (ipMeta * pState) setHyperACKtiveStatus( p * packet_metadata ) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        ps.HyperACKtive = true
        ipMeta.Insert( pKey, ps )
    }
    return ok
}

func (ipMeta * pState) setParentSport( p * packet_metadata, sport int ) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        ps.ParentSport = sport
        ipMeta.Insert( pKey, ps )
    }
    return ok
}

func (ipMeta * pState) getParentSport( p * packet_metadata) int {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        return ps.ParentSport
	}
	return 0
}


func (ipMeta * pState) recordEphemeral( p * packet_metadata, ephemerals []packet_metadata ) bool {

    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
		ps.EphemeralFilters = append( ps.EphemeralFilters, ephemerals... )
        ipMeta.Insert( pKey, ps )
    }
	return ok

}

func (ipMeta * pState) getEphemeralFilters( p * packet_metadata ) ([]packet_metadata, bool) {

    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
		return ps.EphemeralFilters, ok
    }
    return nil, ok

}

func (ipMeta * pState) updateData( p * packet_metadata ) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        ps.Data = true
        ipMeta.Insert( pKey, ps )
    }
    return ok
}

func (ipMeta * pState) getData( p * packet_metadata ) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        return ps.Data
    }
    return false
}



func (ipMeta * pState) getHandshake( p * packet_metadata ) int {
	pKey := constructKey(p)
	ps, ok := ipMeta.Get(pKey)
	if ok {
		return ps.HandshakeNum
	}
	return 0
}

func (ipMeta * pState) incrementCounter( p * packet_metadata ) bool {

	pKey := constructKey(p)
	ps, ok := ipMeta.Get(pKey)
	if !ok {
		return false
	}
	ps.Packet.incrementCounter()
	ipMeta.Insert( pKey, ps )
	return true

}


func (ipMeta * pState) remove( packet *packet_metadata ) *packet_metadata {
	packet.ACKed = ipMeta.getAck( packet )
	packetKey := constructKey(packet)
	ipMeta.Remove( packetKey )
	return packet
}

func verifySA( pMap *packet_metadata, pRecv *packet_metadata ) bool {

	if pRecv.SYN && pRecv.ACK {
		if ( pRecv.Acknum == pMap.Seqnum + 1 ) {
			return true
		}
	} else {

		if ((pRecv.Seqnum == ( pMap.Seqnum )) || (pRecv.Seqnum == ( pMap.Seqnum + 1 ))) {
			if ( pRecv.Acknum == ( pMap.Acknum + pMap.LZRResponseL ) ) {
				return true
			}
			if pRecv.Acknum == 0 { //for RSTs
				return true
			}
		}
	}
	return false

}

//TODO: eventually remove the act of updating packet with hyperactive flag to 
// another packet func
func ( ipMeta * pState ) verifyScanningIP( pRecv *packet_metadata ) bool {

	pRecvKey := constructKey(pRecv)
	//first check that IP itself is being scanned
	ps, ok := ipMeta.Get(pRecvKey)
	if !ok {
		return false
	}
	pMap := ps.Packet

	//second check that 4-tuple matches with default packet
	if (( pMap.Saddr == pRecv.Saddr ) && (pMap.Dport == pRecv.Dport) &&
		(pMap.Sport == pRecv.Sport) ) {

		if verifySA( pMap, pRecv) {
			return true
		}
	}

	/*//lets re-query for the ACKtive packets
	pRecv.HyperACKtive = true
	pRecvKey = constructKey(pRecv)
	ps, ok = ipMeta.Get( pRecvKey )
	if !ok {
		pRecv.HyperACKtive = false
		return false
	}
	pMap = ps.Packet

	if verifySA( pMap, pRecv) {
		return true
	}
	pRecv.HyperACKtive = false
	*/
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


