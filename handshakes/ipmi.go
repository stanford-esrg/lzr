package handshakes

import "github.com/stanford-esrg/lzr/handshakes/ipmi"

func init() {
	ipmi.RegisterHandshake()
}

