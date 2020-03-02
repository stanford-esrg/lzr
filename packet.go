package main

import (
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "log"
	"net"
	"bytes"
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


func getPacketMetadata( ip *layers.IPv4, tcp *layers.TCP ) packet_metadata {

	var packet packet_metadata
	packet.Saddr = ip.SrcIP.String()
	packet.Daddr = ip.DstIP.String()
	packet.Sport = int(tcp.SrcPort)
	packet.Dport = int(tcp.DstPort)
	packet.Seqnum = int(tcp.Seq)
	packet.Acknum = int(tcp.Ack)
	packet.Window = int(tcp.Window)

	return packet
}

func (synack *packet_metadata) windowZero() bool {
    if synack.Window == 0 {
        return true
    }
    return false
}


/*  Packet Ops */

func getSourceMacAddr() (addr net.HardwareAddr) {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				// Don't use random as we have a real address
				addr = i.HardwareAddr
				//break
			}
		}
	}
	return
}

func getHostMacAddr() (addr net.HardwareAddr) {
	return net.HardwareAddr{0xa8, 0x1e, 0x84, 0xce, 0x64, 0x5f}
}

func constructEthLayer() (eth *layers.Ethernet) {

    ethernetLayer := &layers.Ethernet{
        SrcMAC: getSourceMacAddr(),
        DstMAC: getHostMacAddr(),
        //EthernetType: layers.EthernetTypeARP,
        EthernetType: layers.EthernetTypeIPv4,
    }

	return ethernetLayer

}


func constructAck( synack packet_metadata ) []byte {

    //data := []byte("\n")
    data := getData(string(synack.Saddr))

	ethernetLayer := constructEthLayer()

    ipLayer := &layers.IPv4{
        SrcIP: net.ParseIP(synack.Daddr),
        DstIP: net.ParseIP(synack.Saddr),
    TTL : 64,
    Protocol: layers.IPProtocolTCP,
    Version: 4,
    }

    tcpLayer := &layers.TCP{
        SrcPort: layers.TCPPort(synack.Dport),
        DstPort: layers.TCPPort(synack.Sport),
    Seq: uint32(synack.Acknum),
    Ack: uint32(synack.Seqnum+1),
    Window: 8192,
    ACK: true,
    }

    buffer = gopacket.NewSerializeBuffer()
    options := gopacket.SerializeOptions{
        ComputeChecksums: true,
        FixLengths:       true,
    }
    tcpLayer.SetNetworkLayerForChecksum(ipLayer)
    // And create the packet with the layers
    if err := gopacket.SerializeLayers(buffer, options,
		ethernetLayer,
        ipLayer,
        tcpLayer,
        gopacket.Payload(data),
    ); err != nil {
        log.Fatal(err)

	}
    outPacket := buffer.Bytes()
    return outPacket

}


func constructRST( ack packet_metadata ) []byte {

	ethernetLayer := constructEthLayer()

    ipLayer := &layers.IPv4{
        SrcIP: net.ParseIP(ack.Daddr),
        DstIP: net.ParseIP(ack.Saddr),
    TTL : 64,
    Protocol: layers.IPProtocolTCP,
    Version: 4,
    }

    tcpLayer := &layers.TCP{
        SrcPort: layers.TCPPort(ack.Dport),
        DstPort: layers.TCPPort(ack.Sport),
    Seq: uint32(ack.Acknum), //NOT SURE
    Ack: 0,
    Window: 0,
    RST: true,
    }

    buffer = gopacket.NewSerializeBuffer()
    options := gopacket.SerializeOptions{
        ComputeChecksums: true,
        FixLengths:       true,
    }
    tcpLayer.SetNetworkLayerForChecksum(ipLayer)
    // And create the packet with the layers
    if err := gopacket.SerializeLayers(buffer, options,
		ethernetLayer,
        ipLayer,
        tcpLayer,
    ); err != nil {
        log.Fatal(err)

    }
    outPacket := buffer.Bytes()
    return outPacket


}



func verifyScanningIP( pRecv packet_metadata, ipMeta * map[string]packet_metadata  ) bool {


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

