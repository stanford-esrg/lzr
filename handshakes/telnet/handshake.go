package telnet

import (
	"github.com/stanford-esrg/lzr"
	"strings"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

//server first protocol
func (h *HandshakeMod) GetData( dst string ) []byte {
	data := []byte("")
	return data
}

func (h *HandshakeMod) Verify( data string ) string {
	if len(data) < 2 {
		return ""
	}
	if strings.Contains( ToLower(data), "telnet" ) {
		return "telnet"
	}

	if (  data[1] == byte(0xff) || data[1] == byte(0xfe) ||
			data[1] == byte(0xfd) || data[1] == byte(0xfc) ||
				data[1] == byte(0xfb) ) {
		return "telnet"
	}
	//response to HTTP
	/*if strings.Contains( data, "Auth Result") &&
		strings.Contains( data, "Server Time") &&
		strings.Contains( data, "IP Address") {
		return "telnet"
	}*/

	return ""
}

//more efficient string toLower
//https://github.com/golang/go/issues/17859

func ToLower(s string) string {
                b := make([]byte, len(s))
                for i, c := range s {
                        if c >= 'A' && c <= 'Z' {
                                c += 32
                        }
                        b[i] = byte(c)
                }
                return string(b)
}


func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "telnet", &h )
}

