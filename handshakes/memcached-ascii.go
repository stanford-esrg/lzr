package handshakes

import "github.com/stanford-esrg/lzr/handshakes/memcached_ascii"

func init() {
	memcached_ascii.RegisterHandshake()
}

