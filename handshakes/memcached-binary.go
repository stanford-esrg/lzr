package handshakes

import "lzr/handshakes/memcached_binary"

func init() {
	memcached_binary.RegisterHandshake()
}

