package lzr

import (
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "time"
	"encoding/json"
    "log"
	"math"
	//"fmt"
)

var (

	ACK			string = "ack"
	SYN_ACK		string = "sa"
	DATA		string = "data"

)

type packet_state struct {
	HandshakeNum	int
	Packet			*packet_metadata
}

type packet_metadata struct {

	Saddr		    string		`json:"saddr"`
	Daddr		    string		`json:"daddr"`
	Sport		    int			`json:"sport"`
	Dport		    int			`json:"dport"`
	Seqnum		    int			`json:"seqnum"`
	Acknum		    int			`json:"acknum"`
	Window		    int			`json:"window"`
	Counter		    int

    ACK             bool
    SYN             bool
    RST             bool
    FIN             bool
    PUSH            bool
    ValFail         bool

    HandshakeNum	int			`json:"-"`
    Fingerprint     string		`json:"fingerprint,omitempty"`
	Timestamp	    time.Time
    LZRResponseL    int			`json:"-"`
	ExpectedRToLZR  string		`json:"expectedRToLZR,omitempty"`
    Data            string		`json:"data,omitempty"`
    Processing      bool		`json:"-"`
}


func ReadLayers( ip *layers.IPv4, tcp *layers.TCP ) *packet_metadata {

    packet := &packet_metadata{
        Saddr: ip.SrcIP.String(),
        Daddr: ip.DstIP.String(),
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


func convertToPacketM ( packet *gopacket.Packet ) ( *packet_metadata ) {
        tcpLayer := (*packet).Layer(layers.LayerTypeTCP)
        if tcpLayer != nil {
            tcp, _ := tcpLayer.(*layers.TCP)
            ipLayer := (*packet).Layer(layers.LayerTypeIPv4)
            ip, _ := ipLayer.(*layers.IPv4)
            metapacket := ReadLayers(ip,tcp)
            return metapacket
        }
        return nil
}

func convertToPacket ( input string ) *packet_metadata  {


        synack := &packet_metadata{}
        //expecting ip,sequence number, acknumber,windowsize
        err := json.Unmarshal( []byte(input),synack )
		synack.Processing = true
        if err != nil {
            log.Fatal(err)
            return nil
        }
        return synack
}


func (packet * packet_metadata) updatePacketFlow()  {

	//creating a new sourceport to send from 
	//and incrementing the handshake we are trying

	newsrcprt := math.Mod(float64(packet.Dport),65535)+1
	packet.Dport = int(newsrcprt)
	packet.HandshakeNum += 1
	packet.ExpectedRToLZR = SYN_ACK
	packet.Seqnum = packet.Acknum
	packet.Acknum = 0

}

func (packet * packet_metadata) windowZero() bool {
    if packet.Window == 0 {
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

func (packet * packet_metadata) fingerprintData() {

    //return fingerprintResponse( payload )
    packet.Fingerprint = fingerprintResponse( packet.Data )

}

