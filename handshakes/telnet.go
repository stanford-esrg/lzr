package handshakes

import "lzr/handshakes/telnet"

func init() {
	telnet.RegisterHandshake()
}

