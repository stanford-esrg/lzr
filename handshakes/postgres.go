package handshakes

import "lzr/handshakes/postgres"

func init() {
	postgres.RegisterHandshake()
}

