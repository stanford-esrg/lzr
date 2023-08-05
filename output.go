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
	HyperACKtive	int
}


func Summarize( t time.Duration ) {
	fmt.Fprintln(os.Stderr, "Runtime:", t)
	out, _ := json.Marshal( summaryLZR )
	fmt.Fprintln(os.Stderr, string(out))
	//print out fingerprints
	for k, v := range GetFingerprints() {
		fmt.Fprintln(os.Stderr, k +":", v)
	}
}

func addToSummary( packet *packet_metadata ) {

	summaryLZR.TotalResponses  += 1

	if packet.HyperACKtive {
		summaryLZR.HyperACKtive +=1
		return
	}
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
	if  !packet.SYN	&& packet.ACK {
		summaryLZR.Resp_ack += 1
	}
}

func ( f *output_file ) Record( packet *packet_metadata, handshakes []string ) {

	packet.fingerprintData()

	if FeedZGrab() {
		if packet.Fingerprint != "" {
			fmt.Println( packet.Saddr + ", ," + packet.Fingerprint )
		}
	}

	addToSummary( packet )

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

	var f *os.File
	if fname == "-" {
		f = os.Stdout
	} else {
		file, err := os.OpenFile( fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777 )
		if err != nil {
			panic(err)
			log.Fatal(err)
		}
		f = file
	}

	o := &output_file{
		F: bufio.NewWriter(f),
	}

	return o
}

