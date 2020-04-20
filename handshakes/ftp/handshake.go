package ftp

import (
	"lzr"
	"unicode"
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
	if strings.Contains( ToLower(data), "ftp" ) {
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
	if s == "" { // quick return for empty strings
		return s
	}
	if isASCII(s) { // optimize for ascii condition
		b := make([]byte, len(s))
		for i, c := range s {
			if c >= 'A' && c <= 'Z' {
				c += 32
			}
			b[i] = byte(c)
		}
		return string(b)
	}
	//should i even call ToLower for non-ascii??
	return strings.Map(unicode.ToLower, s)
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			return false
		}
	}
	return true
}

