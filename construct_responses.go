package lzr

import (
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "log"
	"net"
	"bytes"
    //"fmt"
)


/*  Packet Ops */
var (
	source_mac = getSourceMacAddr()
)
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
        SrcMAC: source_mac,
        DstMAC: getHostMacAddr(),
        //EthernetType: layers.EthernetTypeARP,
        EthernetType: layers.EthernetTypeIPv4,
    }

	return ethernetLayer

}

func constructSYN( p *packet_metadata ) []byte {

	ethernetLayer := constructEthLayer()

	ipLayer := &layers.IPv4{
        SrcIP: net.ParseIP(p.Daddr),
        DstIP: net.ParseIP(p.Saddr),
		TTL : 64,
		Protocol: layers.IPProtocolTCP,
		Version: 4,
    }

	tcpLayer := &layers.TCP{
		//change srcport slightly
        SrcPort: layers.TCPPort(p.Dport),
        DstPort: layers.TCPPort(p.Sport),
		Seq: uint32(p.Seqnum),
		Ack: uint32(p.Acknum),
		Window: 65535,
		SYN: true,
    }

    buffer := gopacket.NewSerializeBuffer()
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



func constructData( handshake Handshake, p *packet_metadata, ack bool, push bool) ([]byte, []byte) {

    //data := []byte("\n")

    data := handshake.GetData( string(p.Saddr) )

	ethernetLayer := constructEthLayer()

    ipLayer := &layers.IPv4{
        SrcIP: net.ParseIP(p.Daddr),
        DstIP: net.ParseIP(p.Saddr),
		TTL : 64,
		Protocol: layers.IPProtocolTCP,
		Version: 4,
    }

    tcpLayer := &layers.TCP{
        SrcPort: layers.TCPPort(p.Dport),
        DstPort: layers.TCPPort(p.Sport),
		Seq: uint32(p.Acknum),
		Ack: uint32(p.Seqnum+1),
		Window: 65535,
		ACK: ack,
		PSH: push,
    }

    buffer := gopacket.NewSerializeBuffer()
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
    return outPacket,data

}


func constructRST( ack *packet_metadata ) []byte {

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

    buffer := gopacket.NewSerializeBuffer()
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



