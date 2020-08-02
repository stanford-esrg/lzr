package handshakes

import "github.com/stanford-esrg/lzr/handshakes/modbus"

func init() {
	modbus.RegisterHandshake()
}

