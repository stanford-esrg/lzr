package handshakes

import "github.com/stanford-esrg/lzr/handshakes/memcached_binary"

func init() {
	memcached_binary.RegisterHandshake()
}

