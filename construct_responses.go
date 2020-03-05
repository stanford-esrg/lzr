package main

import (
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "log"
	"net"
	"bytes"
)


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


func constructData( synack packet_metadata, ack bool, push bool) []byte {

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
    ACK: ack,
    PSH: push,
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



