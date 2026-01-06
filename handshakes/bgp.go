package handshakes

import "github.com/stanford-esrg/lzr/handshakes/bgp"

func init() {
	bgp.RegisterHandshake()
}
