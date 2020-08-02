package handshakes

import "github.com/stanford-esrg/lzr/handshakes/amqp"

func init() {
	amqp.RegisterHandshake()
}

