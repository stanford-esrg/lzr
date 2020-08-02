package smtp

import (
	"github.com/stanford-esrg/lzr"
	"strings"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {
	//send EHLO
	data := []byte("\x45\x48\x4c\x4f\x0d\x0a")
    return data
}

func (h *HandshakeMod) Verify( data string ) string {

    dl :=  ToLower(data)
	if strings.Contains( dl, "smtp" ) || strings.Contains( dl, "ehlo" ) {
         return "smtp"
	}
    return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "smtp", &h )
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
