package tls

import (
	"strings"
	"bytes"
	"github.com/stanford-esrg/lzr"
	"math/rand"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {

	data := []byte("\x16\x03\x01\x00\x75\x01\x00\x00\x71\x03\x03")
	token := make([]byte, 32)
	rand.Read(token)
	data2 := []byte("\x00\x00\x1a\xc0\x2f\xc0\x2b\xc0\x11\xc0\x07\xc0\x13\xc0\x09\xc0\x14\xc0\x0a\x00\x05\x00\x2f\x00\x35\xc0\x12\x00\x0a\x01\x00\x00\x2e\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\x0a\x00\x08\x00\x06\x00\x17\x00\x18\x00\x19\x00\x0b\x00\x02\x01\x00\x00\x0d\x00\x0a\x00\x08\x04\x01\x04\x03\x02\x01\x02\x03\xff\x01\x00\x01\x00")
	data = append(data, token...)
	data = append(data, data2...)
	return data
}

func (h *HandshakeMod) Verify( data string ) string {

	datab := []byte(data)
	if len(datab) < 3 {
		return ""
	}
	if strings.Contains( data, "HTTPS" ) {
		return "tls"
	}

	//http://blog.fourthbit.com/2014/12/23/traffic-analysis-of-an-ssl-slash-tls-session/
	// Record Type Values       dec      hex
	// -------------------------------------
	// CHANGE_CIPHER_SPEC        20     0x14
	// ALERT                     21     0x15
	// HANDSHAKE                 22     0x16
	// APPLICATION_DATA          23     0x17
	//Version Values            dec     hex
	// -------------------------------------
	// SSL 3.0                   3,0  0x0300
	// TLS 1.0                   3,1  0x0301
	// TLS 1.1                   3,2  0x0302
	// TLS 1.2                   3,3  0x0303
	// TLS 1.3                   3,4  0x0304

	if bytes.Equal([]byte(data[0:1]), []byte{0x16} ) ||
		bytes.Equal([]byte(data[0:1]), []byte{0x14} ) ||
		bytes.Equal([]byte(data[0:1]), []byte{0x15} ) ||
		bytes.Equal([]byte(data[0:1]), []byte{0x17} )  {
		if bytes.Equal([]byte(data[1:3]),[]byte{0x03,0x01} ) ||
			bytes.Equal([]byte(data[1:3]),[]byte{0x03,0x02} ) ||
			bytes.Equal([]byte(data[1:3]),[]byte{0x03,0x03} ) ||
			bytes.Equal([]byte(data[1:3]),[]byte{0x03,0x04} ) {
			return "tls"
		} else if bytes.Equal([]byte(data[1:3]),[]byte{0x03,0x00} ) {
			return "ssl"
		}
	}
	return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "tls", &h )
}

