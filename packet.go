package main

import (
    "github.com/google/gopacket/layers"
    "time"
)

var (

	ACK			string = "ack"
	SYN_ACK		string = "sa"
	DATA		string = "data"

)

type packet_metadata struct {

	Saddr		string	`json:"saddr"`
	Daddr		string	`json:"daddr"`
	Sport		int		`json:"sport"`
	Dport		int		`json:"dport"`
	Seqnum		int		`json:"seqnum"`
	Acknum		int		`json:"acknum"`
	Window		int		`json:"window"`
	Counter		int
    ACK         bool
    SYN         bool
    RST         bool
    FIN         bool
    PUSH        bool

    Fingerprint string
	Timestamp	time.Time
    ResponseL   int
	ExpectedR	string
	//SourceQ     string  //might not need this
    Data        string
    //Closed      bool    //might not need: 
                        //to see if connection has been closed in general
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
        ACK: tcp.ACK,
        SYN: tcp.SYN,
        RST: tcp.RST,
        FIN: tcp.FIN,
        PUSH: tcp.PSH,
        Data: string(tcp.Payload),
        Timestamp: time.Now(),
        ExpectedR: "",
        Counter: 0,
    }
	return packet
}

func (synack *packet_metadata) windowZero() bool {
    if synack.Window == 0 {
        return true
    }
    return false
}


func (packet * packet_metadata) updateResponse( state string ) {

	packet.ExpectedR = state

}

func (packet * packet_metadata) updateResponseL( data []byte ) {

	packet.ResponseL = len( data )

}
func (packet * packet_metadata) incrementCounter() {

    packet.Counter += 1

}

func (packet * packet_metadata) updateTimestamp() {

	packet.Timestamp = time.Now()

}

func (packet * packet_metadata) updateData( payload string ) {

	packet.Data = payload

}

func (packet * packet_metadata) fingerprintData() {

    //return fingerprintResponse( payload )
    packet.Fingerprint = fingerprintResponse( packet.Data )

}

