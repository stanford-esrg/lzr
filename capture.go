package main

import (
    "fmt"
    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
    "github.com/google/gopacket/layers"
    "log"
    "time"
    "net/http"
    "net/http/httputil"
    //"net"
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


func constructAck( ip *layers.IPv4, tcp *layers.TCP, ethernet *layers.Ethernet ) []byte {


    //data := []byte("\n")
    data := getData(string(ip.SrcIP))

    ipLayer := &layers.IPv4{
        SrcIP: ip.DstIP,
        DstIP: ip.SrcIP,
	TTL : 64,
	Protocol: layers.IPProtocolTCP,
	Version: 4,
    }

    ethernetLayer := &layers.Ethernet{
        SrcMAC: ethernet.DstMAC,
        DstMAC: ethernet.SrcMAC,
	EthernetType: layers.EthernetTypeIPv4,
    }

    //SEq and Ack not working
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
        ethernetLayer,
        ipLayer,
        tcpLayer,
        gopacket.Payload(data),
    ); err != nil {
        log.Fatal(err)
    }

    fmt.Println(buffer)

    outPacket := buffer.Bytes()
    //packet := gopacket.NewPacket(outPacket, layers.LayerTypeEthernet, gopacket.Default)
    return outPacket

}


func main() {
    // Open device
    handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
    if err != nil {
        log.Fatal(err)
    }
    defer handle.Close()

    // Use the handle as a packet source to process all packets
    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    for packet := range packetSource.Packets() {
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
	    tcp, _ := tcpLayer.(*layers.TCP)
	    ipLayer := packet.Layer(layers.LayerTypeIPv4)
            ip, _ := ipLayer.(*layers.IPv4)
	    ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	    ethernet, _ := ethernetLayer.(*layers.Ethernet)

	    //for every s/a send an ack with data
	    if tcp.SYN && tcp.ACK {
		outPacket := constructAck(ip,tcp,ethernet)
                err = handle.WritePacketData(outPacket)
                if err != nil {
                    log.Fatal(err)
                }
            }

	    //for every ack received, mark as accepting data
	    if !tcp.SYN && tcp.ACK {

	       

	    }

        }
    }
}
