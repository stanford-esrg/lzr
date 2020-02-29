package main

import (
    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
    "github.com/google/gopacket/layers"
    "log"
    "time"
    "net/http"
    "net/http/httputil"
    "io"
    "bufio"
    "os"
    "fmt"
	"encoding/json"
	"net"
	"bytes"
)

var (
    device       string = "ens8"
    snapshot_len int32  = 1024
    promiscuous  bool   = false
    err          error
    timeout      time.Duration = 5 * time.Second
    handle       *pcap.Handle
    buffer       gopacket.SerializeBuffer
)

type packet_metadata struct {

	Saddr		string	`json:"saddr"`
	Daddr		string	`json:"daddr"`
	Sport		int		`json:"sport"`
	Dport		int		`json:"dport"`
	Seqnum		int		`json:"seqnum"`
	Acknum		int		`json:"acknum"`
	Window		int		`json:"window"`
}


/* FUNCS */

func getData( dst string ) []byte {

        req, _ := http.NewRequest("GET","/",nil)
        req.Host =  dst
        req.Header.Add("Host",dst)
        req.Header.Set("User-Agent","Mozilla/5.0 zgrab/0.x")
        req.Header.Set("Accept","*/*")
        req.Header.Set("Accept-Encoding","gzip")
        data, _ := httputil.DumpRequest(req, false)
	return data

}

//TODO: replace with making layers into synack packet_metadata and then pass into constructAck
func constructAckFromStream( ip *layers.IPv4, tcp *layers.TCP, ethernet *layers.Ethernet ) []byte {


    //data := []byte("\n")
    data := getData(string(ip.SrcIP))

    ipLayer := &layers.IPv4{
        SrcIP: ip.DstIP,
        DstIP: ip.SrcIP,
	TTL : 64,
	Protocol: layers.IPProtocolTCP,
	Version: 4,
    }

    /*ethernetLayer := &layers.Ethernet{
        SrcMAC: ethernet.DstMAC,
        DstMAC: ethernet.SrcMAC,
	EthernetType: layers.EthernetTypeIPv4,
    }*/

    tcpLayer := &layers.TCP{
        SrcPort: tcp.DstPort,
        DstPort: tcp.SrcPort,
	Seq: tcp.Ack,
	Ack: tcp.Seq+1,
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
    //    ethernetLayer,
        ipLayer,
        tcpLayer,
        gopacket.Payload(data),
    ); err != nil {
        log.Fatal(err)
    }

    outPacket := buffer.Bytes()
    return outPacket

}

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


func  windowZero(synack packet_metadata) bool {


	if synack.Window == 0 {
		return true
	}
	return false

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



func ackZMap(input string) {

        fmt.Println(input)

        var synack packet_metadata
        //expecting ip,sequence number, acknumber,windowsize
        err = json.Unmarshal( []byte(input),&synack)
        if err != nil {
            log.Fatal(err)
        }

        //TODO: check that ip_metadata contains what we want (saddr,seq,ack,window)

        if windowZero(synack) {
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
		return

}

func handlePcap(packet gopacket.Packet) {


        tcpLayer := packet.Layer(layers.LayerTypeTCP)
        if tcpLayer != nil {
            tcp, _ := tcpLayer.(*layers.TCP)
            ipLayer := packet.Layer(layers.LayerTypeIPv4)
            ip, _ := ipLayer.(*layers.IPv4)

			fmt.Println(packet)
            packet := getPacketMetadata(ip,tcp)

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


func main() {

    //read in config 
    //port := parse()
    //fmt.Println("%s",port)

    //read in s/a sent in by zmap
	_, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	//routine to read in from ZMap
	zmapIncoming := make(chan string)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {

			//Read from ZMap
			input, err := reader.ReadString(byte('\n'))
			if err != nil && err == io.EOF {
				return
			}
			zmapIncoming <- input
		}

	}()


	//routine to read in from pcap
	pcapIncoming := make(chan gopacket.Packet)
	go func() {
		// Open device
		handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, 0) //timeout
		if err != nil {
			log.Fatal(err)
		}
		defer handle.Close()
		// Use the handle as a packet source to process all packets
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

		for {

			//Read from pcap
			packet, err := packetSource.NextPacket()
			if err == io.EOF {
				return
			} else if err != nil {
				log.Println("Error:", err)
				continue
			}
			pcapIncoming <- packet
		}

	}()

	//read from both zmap and pcap
	for {
		select {
			case input := <-zmapIncoming:
				ackZMap(input)
			case input := <-pcapIncoming:
				handlePcap(input)
		}
	}

} //end of main
