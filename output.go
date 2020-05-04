package lzr

import (
    "encoding/json"
    "log"
    "os"
	"bufio"
	"time"
    "fmt"
)

var (
	summaryLZR = &summary{}
)

type output_file struct {

	F	 *bufio.Writer
}

type summary struct {

	TotalResponses	int
	ZeroWindow		int
	ACKed			int
	Data			int
	No_SYNACK		int
	Rst				int
	Fin				int
	Resp_ack		int
	Resp_push		int
}


func Summarize( t time.Duration ) {
	fmt.Println( "Runtime:",t )
	out, _ := json.Marshal( summaryLZR )
	fmt.Println( string(out) )
	//print out fingerprints
	for k, v := range GetFingerprints() {
		fmt.Println(k +":", v)
	}
}

func addToSummary( packet *packet_metadata ) {

	summaryLZR.TotalResponses  += 1
	if packet.Window == 0 {
		summaryLZR.ZeroWindow += 1
	}
	if packet.ACKed {
		summaryLZR.ACKed += 1
	}
    if packet.RST {
        summaryLZR.Rst += 1
    }
    if packet.FIN {
        summaryLZR.Fin += 1
    }
	if packet.ExpectedRToLZR == SYN_ACK {
		summaryLZR.No_SYNACK += 1
	}
	if packet.Data != "" {
		summaryLZR.Data += 1
	}
	if packet.Counter == 0  && packet.ACKed {
		summaryLZR.Resp_ack += 1
	}
	if packet.Counter == 1 && packet.ACKed {
		summaryLZR.Resp_push += 1
	}
}

func ( f *output_file ) Record( packet packet_metadata, handshakes []string ) {

	/*if packet.Data != "" {
		fmt.Println( packet.Saddr + ", , " + handshakes[ packet.HandshakeNum ] )
	}*/

    packet.fingerprintData()
	addToSummary( &packet )

    out, _ := json.Marshal( packet )
    _,err := (f.F).WriteString( string(out) )
    if err != nil {
        f.F.Flush()
        panic(err)
		log.Fatal(err)
	}
    _,err = (f.F).WriteString( "\n" )
    if err != nil {
        f.F.Flush()
        panic(err)
		log.Fatal(err)
	}
    return
}


func InitFile( fname string ) *output_file {

    f, err := os.OpenFile( fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777 )

    if err != nil {
        panic(err)
		log.Fatal(err)
    }

    o := &output_file{
		F: bufio.NewWriter(f),
    }

    return o
}

