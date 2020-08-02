package ftp

import (
	"github.com/stanford-esrg/lzr"
	//"unicode"
	"strings"
	"unicode/utf8"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {
	data := []byte("") // just wait for banner
    return data
}

func (h *HandshakeMod) Verify( data string ) string {
	if data == "" || !isASCII(data) {
		return ""
	} else if strings.Contains( ToLower(data), "ftp" ) {
		return "ftp"
	}
	if len(data) < 4 {
		return ""
	}
	if strings.Contains( data[0:3], "220" ) ||
		strings.Contains( data[0:3], "421" ) ||
		strings.Contains( data[0:3], "530" ) ||
		strings.Contains( data[0:3], "550" ) ||
		strings.Contains( data[0:3], "230" ) {
		return "ftp"
	}

    return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "ftp", &h )
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

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			return false
		}
	}
	return true
}

