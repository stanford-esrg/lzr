package handshakes

import "lzr/handshakes/memcached_ascii"

func init() {
	memcached_ascii.RegisterHandshake()
}

