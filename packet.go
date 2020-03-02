package main

import (
    "github.com/google/gopacket/layers"
)



type packet_metadata struct {

	Saddr		string	`json:"saddr"`
	Daddr		string	`json:"daddr"`
	Sport		int		`json:"sport"`
	Dport		int		`json:"dport"`
	Seqnum		int		`json:"seqnum"`
	Acknum		int		`json:"acknum"`
	Window		int		`json:"window"`
	State		string

}


func NewPacket( ip *layers.IPv4, tcp *layers.TCP ) *packet_metadata {

    packet := &packet_metadata{
        Saddr: ip.SrcIP.String(),
        Daddr: ip.DstIP.String(),
        Sport: int(tcp.SrcPort),
        Dport: int(tcp.DstPort),
        Seqnum: int(tcp.Seq),
        Acknum: int(tcp.Ack),
        Window: int(tcp.Window),
        State: "",
    }
	return packet
}

func (synack *packet_metadata) windowZero() bool {
    if synack.Window == 0 {
        return true
    }
    return false
}


func ( pRecv *packet_metadata ) verifyScanningIP(ipMeta * map[string]packet_metadata  ) bool {


	//first check that IP itself is being scanned
	pSent, ok := (*ipMeta)[pRecv.Saddr]
	if !ok {
		return false
	}
	//second check that 4-tuple matches
	if (( pSent.Saddr == pRecv.Daddr ) && (pSent.Dport == pRecv.Sport) &&
		(pSent.Sport == pRecv.Dport)) {
		return true
	}
	//TODO: check seq & ack and check state that we expect(?)

	return false

}

