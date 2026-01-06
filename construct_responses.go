/*
Copyright 2020 The Board of Trustees of The Leland Stanford Junior University

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package lzr

import (
	"bytes"
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	//"fmt"
)

/*  Packet Ops */

func saveHostMacAddr(packet *packet_metadata) {
	dest_mac = packet.getSourceMac()
}

func getSourceMacAddr() (addr string) {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				if i.Name != getDevice() {
					continue
				}
				addr = i.HardwareAddr.String()
			}
		}
	}
	return addr
}

func constructEthLayer() (eth *layers.Ethernet) {

	smac, _ := net.ParseMAC(source_mac)
	dmac, _ := net.ParseMAC(dest_mac)

	ethernetLayer := &layers.Ethernet{
		SrcMAC: smac,
		DstMAC: dmac,
		//EthernetType: layers.EthernetTypeARP,
		EthernetType: layers.EthernetTypeIPv4,
	}
	if IPv6Enabled() {
		ethernetLayer.EthernetType = layers.EthernetTypeIPv6
	}

	return ethernetLayer

}

/* NOTE: constructing RESPONSE SYN.
 * so Daddr/Saddr etc will be inverted in the process
 */
func constructSYN(p *packet_metadata) []byte {

	ethernetLayer := constructEthLayer()

	var ipLayer4 *layers.IPv4
	var ipLayer6 *layers.IPv6
	if IPv6Enabled() {
		ipLayer6 = &layers.IPv6{
			SrcIP:      net.ParseIP(p.Daddr),
			DstIP:      net.ParseIP(p.Saddr),
			Version:    6,
			HopLimit:   64,
			NextHeader: layers.IPProtocolTCP,
		}
	} else {
		ipLayer4 = &layers.IPv4{
			SrcIP:    net.ParseIP(p.Daddr),
			DstIP:    net.ParseIP(p.Saddr),
			TTL:      64,
			Protocol: layers.IPProtocolTCP,
			Version:  4,
		}
	}

	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(p.Dport),
		DstPort: layers.TCPPort(p.Sport),
		Seq:     uint32(p.Seqnum),
		Ack:     uint32(p.Acknum),
		Window:  uint16(p.Window), //65535,
		SYN:     true,
	}

	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}
	if IPv6Enabled() {
		tcpLayer.SetNetworkLayerForChecksum(ipLayer6)
		// And create the packet with the layers
		if err := gopacket.SerializeLayers(buffer, options,
			ethernetLayer,
			ipLayer6,
			tcpLayer,
		); err != nil {
			log.Fatal(err)

		}
	} else {
		tcpLayer.SetNetworkLayerForChecksum(ipLayer4)
		// And create the packet with the layers
		if err := gopacket.SerializeLayers(buffer, options,
			ethernetLayer,
			ipLayer4,
			tcpLayer,
		); err != nil {
			log.Fatal(err)

		}
	}

	outPacket := buffer.Bytes()
	return outPacket
}

/* NOTE: constructing RESPONSE.
 * so Daddr/Saddr etc will be inverted in the process
 */
func constructData(handshake Handshake, p *packet_metadata, ack bool, push bool) ([]byte, []byte) {

	//data := []byte("\n")

	data := handshake.GetData(string(p.Saddr))
	if PushDOnly() && !push {
		data = []byte("")
	}
	ethernetLayer := constructEthLayer()

	var ipLayer4 *layers.IPv4
	var ipLayer6 *layers.IPv6
	if IPv6Enabled() {
		ipLayer6 = &layers.IPv6{
			SrcIP:      net.ParseIP(p.Daddr),
			DstIP:      net.ParseIP(p.Saddr),
			Version:    6,
			HopLimit:   64,
			NextHeader: layers.IPProtocolTCP,
		}
	} else {
		ipLayer4 = &layers.IPv4{
			SrcIP:    net.ParseIP(p.Daddr),
			DstIP:    net.ParseIP(p.Saddr),
			TTL:      64,
			Protocol: layers.IPProtocolTCP,
			Version:  4,
		}
	}

	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(p.Dport),
		DstPort: layers.TCPPort(p.Sport),
		Seq:     uint32(p.Acknum),
		Ack:     uint32(p.Seqnum + 1),
		Window:  uint16(p.Window),
		ACK:     ack,
		PSH:     push,
	}

	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}
	if IPv6Enabled() {
		tcpLayer.SetNetworkLayerForChecksum(ipLayer6)
		// And create the packet with the layers
		if err := gopacket.SerializeLayers(buffer, options,
			ethernetLayer,
			ipLayer6,
			tcpLayer,
			gopacket.Payload(data),
		); err != nil {
			log.Fatal(err)
		}
	} else {
		tcpLayer.SetNetworkLayerForChecksum(ipLayer4)
		// And create the packet with the layers
		if err := gopacket.SerializeLayers(buffer, options,
			ethernetLayer,
			ipLayer4,
			tcpLayer,
			gopacket.Payload(data),
		); err != nil {
			log.Fatal(err)
		}
	}
	outPacket := buffer.Bytes()
	return outPacket, data

}

/* NOTE: constructing RESPONSE.
 * so Daddr/Saddr etc will be inverted in the process
 */
func constructRST(ack *packet_metadata) []byte {

	ethernetLayer := constructEthLayer()

	var ipLayer4 *layers.IPv4
	var ipLayer6 *layers.IPv6
	if IPv6Enabled() {
		ipLayer6 = &layers.IPv6{
			SrcIP:      net.ParseIP(ack.Daddr),
			DstIP:      net.ParseIP(ack.Saddr),
			Version:    6,
			HopLimit:   64,
			NextHeader: layers.IPProtocolTCP,
		}
	} else {
		ipLayer4 = &layers.IPv4{
			SrcIP:    net.ParseIP(ack.Daddr),
			DstIP:    net.ParseIP(ack.Saddr),
			TTL:      64,
			Protocol: layers.IPProtocolTCP,
			Version:  4,
		}
	}

	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(ack.Dport),
		DstPort: layers.TCPPort(ack.Sport),
		Seq:     uint32(ack.Acknum), //NOT SURE
		Ack:     0,
		Window:  0,
		RST:     true,
	}

	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}
	if IPv6Enabled() {
		tcpLayer.SetNetworkLayerForChecksum(ipLayer6)
		// And create the packet with the layers
		if err := gopacket.SerializeLayers(buffer, options,
			ethernetLayer,
			ipLayer6,
			tcpLayer,
		); err != nil {
			log.Fatal(err)

		}
	} else {
		tcpLayer.SetNetworkLayerForChecksum(ipLayer4)
		// And create the packet with the layers
		if err := gopacket.SerializeLayers(buffer, options,
			ethernetLayer,
			ipLayer4,
			tcpLayer,
		); err != nil {
			log.Fatal(err)

		}
	}
	outPacket := buffer.Bytes()
	return outPacket

}
