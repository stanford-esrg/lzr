package kubernetes

import (
	"strings"
	//"fmt"
	"github.com/stanford-esrg/lzr"
	"math/rand"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {
	//just tls
	data := []byte("\x16\x03\x01\x00\x75\x01\x00\x00\x71\x03\x03")
	token := make([]byte, 32)
    rand.Read(token)
	data2 := []byte("\x00\x00\x1a\xc0\x2f\xc0\x2b\xc0\x11\xc0\x07\xc0\x13\xc0\x09\xc0\x14\xc0\x0a\x00\x05\x00\x2f\x00\x35\xc0\x12\x00\x0a\x01\x00\x00\x2e\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\x0a\x00\x08\x00\x06\x00\x17\x00\x18\x00\x19\x00\x0b\x00\x02\x01\x00\x00\x0d\x00\x0a\x00\x08\x04\x01\x04\x03\x02\x01\x02\x03\xff\x01\x00\x01\x00")
	data = append(data, token...)
	data = append(data, data2...)
	return data

}

func (h *HandshakeMod) Verify( data string ) string {

	if	 strings.Contains( data, "kubernetes" ) {
         return "kubernetes"
	}
	return ""

}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "kubernetes", &h )
}

