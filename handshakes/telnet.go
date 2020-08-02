package handshakes

import "github.com/stanford-esrg/lzr/handshakes/telnet"

func init() {
	telnet.RegisterHandshake()
}

