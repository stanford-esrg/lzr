package handshakes

import "github.com/stanford-esrg/lzr/handshakes/redis"

func init() {
	redis.RegisterHandshake()
}

