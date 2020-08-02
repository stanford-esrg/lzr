package handshakes

import "github.com/stanford-esrg/lzr/handshakes/tls"

func init() {
	tls.RegisterHandshake()
}

