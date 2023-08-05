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
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"time"
	"encoding/json"
	"log"
	"math"
	"strings"
	"strconv"
	//"fmt"
	//"os"
)

var (

	ACK			string = "ack"
	SYN_ACK		string = "sa"
	DATA		string = "data"

)

type packet_state struct {
	HandshakeNum		int
	Ack					bool
	Data				bool
	HyperACKtive		bool
	EphemeralFilters	[]packet_metadata
	EphemeralRespNum	int
	ParentSport			int			//used for filter packets
	Packet				*packet_metadata
}

type packet_metadata struct {

	Smac				string		`json:"-"`
	Dmac				string		`json:"-"`
	Saddr				string		`json:"saddr"`
	Daddr				string		`json:"daddr"`
	Sport				int			`json:"sport"`
	Dport				int			`json:"dport"`
	Seqnum				int			`json:"seqnum"`
	Acknum				int			`json:"acknum"`
	Window				int			`json:"window"`
	TTL					uint8		`json:"ttl"`
	Counter				int

	ACK					bool
	ACKed				bool
	SYN					bool
	RST					bool
	FIN					bool
	PUSH				bool
	ValFail				bool		`json:"-"`

	HandshakeNum		int
	Fingerprint			string		`json:"fingerprint,omitempty"`
	Timestamp			time.Time
	LZRResponseL		int			`json:"-"`
	ExpectedRToLZR		string		`json:"expectedRToLZR,omitempty"`
	Data				string		`json:"data,omitempty"`
	Processing			bool		`json:"-"`
	HyperACKtive		bool		`json:"ackingFirewall,omitempty"`
}


func ReadLayers( ip *layers.IPv4, tcp *layers.TCP, eth *layers.Ethernet ) *packet_metadata {

	packet := &packet_metadata{
		Smac: eth.SrcMAC.String(),
		Dmac: eth.DstMAC.String(),
		Saddr: ip.SrcIP.String(),
		Daddr: ip.DstIP.String(),
		TTL: ip.TTL,
		Sport: int(tcp.SrcPort),
		Dport: int(tcp.DstPort),
		Seqnum: int(tcp.Seq),
		Acknum: int(tcp.Ack),
		Window: int(tcp.Window),
		ACK: tcp.ACK,
		SYN: tcp.SYN,
		RST: tcp.RST,
		FIN: tcp.FIN,
		PUSH: tcp.PSH,
		Data: string(tcp.Payload),
		Timestamp: time.Now(),
		Counter: 0,
		Processing: true,
		HandshakeNum: 0,
	}
	return packet
}

func convertToPacketM( packet *gopacket.Packet ) *packet_metadata {

	tcpLayer := (*packet).Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		ipLayer := (*packet).Layer(layers.LayerTypeIPv4)
		if ipLayer != nil {
			ip, _ := ipLayer.(*layers.IPv4)

			ethLayer := (*packet).Layer(layers.LayerTypeEthernet)
			if ethLayer != nil {
				eth, _ := ethLayer.(*layers.Ethernet)
				metapacket := ReadLayers(ip,tcp,eth)
				return metapacket
			}
		}
	}
	return nil
}

func convertFromZMapToPacket( input string ) *packet_metadata	{

	synack := &packet_metadata{}
	//expecting ip,sequence number, acknumber,windowsize, sport, dport
	err := json.Unmarshal( []byte(input),synack )
	synack.Processing = true
    synack.SYN = true
    synack.ACK = true
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return synack
}


func convertFromInputListToPacket( input string ) *packet_metadata {

	t := time.Now()
	//expecting ip, port
	input = strings.TrimSuffix(input, "\n")
	s := strings.Split(input,":")
	if len(s) != 2 {
		panic("Error parsing input list")
	}
	saddr, sport_s := s[0], s[1]
	sport, err := strconv.Atoi(sport_s)
	if err != nil {
		panic(err)
		panic("Wrong input list format")
		panic("BAD STUFF IS ABOUT TO HAPPEN")
	}
	if getHostMacAddr() == "" {
		panic("Gateway Mac Address required")
	}
	if getSourceIP() == "" {
		 panic("Source IP required")
	}

	//note that source and dest are inverted
	syn := &packet_metadata{
		Smac: source_mac,
		Dmac: getHostMacAddr(),
        Saddr: saddr,
        Daddr: getSourceIP(),
		Dport: randInt(32768, 61000, t.UnixNano()),
        Sport: sport,
		Seqnum: int(math.Mod(float64(t.UnixNano()),65535)),
        Acknum: 0,
        Window: 65535,
        SYN: true,
        Timestamp: t,
        Counter: 0,
        Processing: true,
        HandshakeNum: 0,
        ExpectedRToLZR: SYN_ACK,
    }

	return syn
}

func randInt(min int, max int, cur int64) int {
	//return min + rand.Intn(max-min)
	return min + int(math.Mod(float64(cur), float64(max-min)))
}

//create a packet to filter out nets like canada
func createFilterPacket( packet *packet_metadata ) *packet_metadata {

	t := time.Now()
	packetFilter := &packet_metadata{
		Smac: packet.Smac,
		Dmac: packet.Dmac,
		Saddr: packet.Saddr,
		Daddr: packet.Daddr,
		Dport: int(math.Mod(float64(packet.Dport),65535)+1),
		Sport: randInt(32768, 61000, t.UnixNano()),
		Seqnum: int(math.Mod(float64(t.UnixNano()),65535)),
		Acknum: 0,
		Window: packet.Window,
		SYN: true,
		Timestamp: t,
		Counter: 0,
		Processing: true,
		HandshakeNum: 0,
		HyperACKtive: true,
		ExpectedRToLZR: SYN_ACK,
	}
	return packetFilter

}


func ( packet * packet_metadata ) updatePacketFlow()  {

	//creating a new sourceport to send from 
	//and incrementing the handshake we are trying

	newsrcprt := math.Mod(float64(packet.Dport),65535)+1
	packet.Dport = int(newsrcprt)
	packet.HandshakeNum += 1
	packet.Counter = 0
	packet.ExpectedRToLZR = SYN_ACK
	packet.Seqnum = packet.Acknum
	packet.Acknum = 0
	packet.Data = ""
	packet.Fingerprint = ""
	packet.SYN = false
	packet.ACK = false
	packet.PUSH = false
	packet.RST = false
	packet.FIN = false
}

func (packet * packet_metadata) windowZero() bool {
	if packet.Window == 0 && packet.SYN && packet.ACK {
		return true
	}
	return false
}

func (packet * packet_metadata) hasData() bool {
	if len(packet.Data) > 0 {
		return true
	}
	return false
}

func (packet * packet_metadata) syncHandshakeNum( h int ) {

	packet.HandshakeNum = h

}

func (packet * packet_metadata) getHandshakeNum() int {
	return packet.HandshakeNum

}

func (packet * packet_metadata) updateResponse( state string ) {

	packet.ExpectedRToLZR = state

}

func (packet * packet_metadata) updateResponseL( data []byte ) {

	packet.LZRResponseL = len( data )

}
func (packet * packet_metadata) incrementCounter() {

	packet.Counter += 1

}

func (packet * packet_metadata) updateTimestamp() {

	packet.Timestamp = time.Now()

}

func (packet * packet_metadata) startProcessing() {

	packet.Processing = true

}

func (packet * packet_metadata) finishedProcessing() {

	packet.Processing = false

}

func (packet * packet_metadata) updateData( payload string ) {

	packet.Data = payload

}


func (packet * packet_metadata) validationFail() {

	packet.ValFail = true

}

func (packet * packet_metadata) getValidationFail() bool {

	return packet.ValFail

}


func (packet * packet_metadata) getSourceMac() string {

    return packet.Smac

}

func (packet * packet_metadata) fingerprintData() {

	packet.Fingerprint = fingerprintResponse( packet.Data )

}

func (packet * packet_metadata) setHyperACKtive( ackingFirewall bool ) {

	packet.HyperACKtive = ackingFirewall

}


