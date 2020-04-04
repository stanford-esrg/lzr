package ipmi

import (
	"bytes"
	//"fmt"
)

func getIPMI() []byte {

   // The generic probe is
    //   06,00,ff,07|00,00 00 00 00,00 00 00 00,09|20,18,c8,81,00,38[8e 04]b5
    //  | RMCP      | IMPI 1.5 header             | IPMI payload
    // v=0, r=0, seq=0, cls=02 = invalid?
    rmcpHeader := &RMCPHeader{
        RMCPVersion1_0,
        0x00,
        0xff,
        MessageClassIPMI,
    }
    ipmiHeader := &IPMISessionHeader_v1_5{
        AuthType:              AuthTypeNone,
        SessionSequenceNumber: 0,
        SessionID:             0,
    }
    command := (&GetChannelAuthenticationCapabilitiesRequest{}).Set(true, 0x0E,  ChannelPrivilegeLevelAdmin)
    // 18 = 0001 1000 -> 00 0110 | 00
    // 81 = 0x40 << 1 + 1
    rqAddr := IPMIAddressByte(0)
    rqAddr.Set(true, 0x40)
    ipmiPayload := &IPMICommandPayload{
        // Responder address.
        RsAddr: IPMIAddressBMC,
        NetFn:  NetFnApp.Pack(0),
        // Chk1 auto-calculated
        RqAddr: rqAddr,

        // Requestor sequence number.
        RqSeq: 0x00,

        // Command byte.
        Cmd: CmdGetChannelAuthenticationCapabilities,
    }
    if err := ipmiPayload.SetData(command); err != nil {
        return []byte("")
    }
    packet, err := ipmiHeader.GetPacket(ipmiPayload)
    if err != nil {
        return []byte("")
    }
    temp := bytes.NewBuffer(nil)
    if _, err := rmcpHeader.Write(temp); err != nil {
        return []byte("")
    }
    if _, err := packet.Write(temp); err != nil {
        return []byte("")
    }

	return temp.Bytes()
}

//note that im basically only verifying rmcp
func verifyIPMI( data string ) string {
	datab := []byte(data)
	var positiveDetectUnknown = []byte{
    0x00, 0x00, 0x00, 0x02, 0x09, 0x00, 0x00, 0x00,
    0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00}


	if bytes.Equal( datab, positiveDetectUnknown) {
		return "ipmi"
	}
	if len(datab) < 4 {
		return ""
	}
	//pg.126
	//https://www.intel.com/content/dam/www/public/us/en/documents/
	//product-briefs/ipmi-second-gen-interface-spec-v2-rev1-1.pdf

	if datab[0] != RMCPVersion1_0 {
		return ""
	}
    if datab[1] != 0x0 { //reserved field
        return ""
    }
	if datab[2] != 0xff { //seq #
        return ""
    }
    if MessageClass(datab[3]).Class() != MessageClassIPMI.Class() {
        return ""
    }
	return "ipmi"
}
