package main

import (
    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
    "github.com/google/gopacket/layers"
    "log"
    "time"
    "fmt"
	"encoding/json"
)



var (
    device       string = "ens8"
    snapshot_len int32  = 1024
    promiscuous  bool   = false
    err          error
    timeout      time.Duration = 5 * time.Second
    handle       *pcap.Handle
    buffer       gopacket.SerializeBuffer
	ACK			string = "ack"
	SYN_ACK		string = "sa"
	DATA		string = "data"
)


func ackZMap(input string, ipMeta * map[string]packet_metadata) {

        fmt.Println(input)

        var synack packet_metadata
        //expecting ip,sequence number, acknumber,windowsize
        err = json.Unmarshal( []byte(input),&synack)
        if err != nil {
            log.Fatal(err)
        }

        //TODO: check that ip_metadata contains what we want (saddr,seq,ack,window)

        if synack.windowZero() {
            //not a real s/a
            return
        }

        //Send Ack with Data
        ack := constructAck(synack)
        fmt.Println("Constructed ack...")
        err = handle.WritePacketData(ack)
        if err != nil {
            log.Fatal(err)
        }

		//add to map
		synack.State = ACK
		(*ipMeta)[synack.Saddr] = synack
		return

}


func handlePcap(packet gopacket.Packet, ipMeta * map[string]packet_metadata) {

        tcpLayer := packet.Layer(layers.LayerTypeTCP)
        if tcpLayer != nil {
            tcp, _ := tcpLayer.(*layers.TCP)
            ipLayer := packet.Layer(layers.LayerTypeIPv4)
            ip, _ := ipLayer.(*layers.IPv4)

            packet := getPacketMetadata(ip,tcp)
			//verify 
			if !verifyScanningIP( packet, ipMeta ) {
				return
			}
			fmt.Println(packet)

            //for every ack received, mark as accepting data
            if (!tcp.SYN) && tcp.ACK {
                //TODO: do something with data
                fmt.Println(tcp.Payload)
                fmt.Println("acked")
                //close connection
                rst := constructRST(packet)
                err = handle.WritePacketData(rst)
            }
        }


}

