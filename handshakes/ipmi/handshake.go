package ipmi

import (
	"bytes"
	"github.com/stanford-esrg/lzr"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {

    //   06,00,ff,07|00,00 00 00 00,00 00 00 00,09|20,18,c8,81,00,38[8e 04]b5
    //  | RMCP      | IMPI 1.5 header             | IPMI payload
    // v=0, r=0, seq=0, cls=02 = invalid?
	data := []byte("\x06\x00\xff\x07\x00\x00\x00\x00\x00\x00\x00\x00\x00\x09\x20\x18\xc8\x81\x00\x38\x8e\x04\xb5")
    return data
}

func (h *HandshakeMod) Verify( data string ) string {

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
	if bytes.Equal([]byte(data[0:4]),[]byte{ 0x06, 0x00, 0xff, 0x07 } ) {
		return "ipmi"
	}
	return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "ipmi", &h )
}

