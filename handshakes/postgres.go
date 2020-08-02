package handshakes

import "github.com/stanford-esrg/lzr/handshakes/postgres"

func init() {
	postgres.RegisterHandshake()
}

