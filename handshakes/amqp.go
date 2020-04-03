package handshakes

import "lzr/handshakes/amqp"

func init() {
	amqp.RegisterHandshake()
}

