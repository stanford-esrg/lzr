package bgp

import (
	"github.com/stanford-esrg/lzr"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

// Send Garbage Data to prompt a notification message
func (h *HandshakeMod) GetData(dst string) []byte {
	data := []byte("\n\n")
	return data
}

func (h *HandshakeMod) Verify(data string) string {
	if len(data) < 16 {
		return ""
	}

	correctBGP := true
	for i := 0; i < 16 && correctBGP; i++ {
		if data[i] != byte(0xff) {
			correctBGP = false
		}
	}

	if correctBGP == true {
		return "bgp"
	}
	return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake("bgp", &h)
}
