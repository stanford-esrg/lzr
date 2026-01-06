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
	if len(data) < 4 {
		return ""
	}
    if strings.Contains( data[0:3], "220" ) ||
		strings.Contains( data[0:3], "421" ) ||
		strings.Contains( data[0:3], "530" ) ||
		strings.Contains( data[0:3], "550" ) ||
		strings.Contains( data[0:3], "230" ) {
		if !strings.Contains( ToLower(data), "ftp" ) {
			// We need an additional check against FTP because it uses the same status codes.
			return "smtp"
		}
		return ""
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
