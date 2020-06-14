package handshakes

import "lzr/handshakes/redis"

func init() {
	redis.RegisterHandshake()
}

