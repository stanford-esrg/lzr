package trueconf

import (
	"strings"

	"github.com/stanford-esrg/lzr"
)

// Handshake implements the lzr.Handshake interface.
type Handshake struct{}

func (h Handshake) GetData(_ string) []byte {
	return []byte("_VS_TRANSPORT_\x00\x00\x99\xef+\x00\x00\x00\x00\x00\x00\x01")
}

func (h Handshake) Verify(data string) string {
	if tail, found := strings.CutPrefix(data, "_VS_TRANSPORT_"); found {
		if strings.Contains(tail, "#vcs") {
			return "trueconf"
		}
	}
	return ""
}

func RegisterHandshake() {
	var h Handshake
	lzr.AddHandshake("trueconf", &h)
}
