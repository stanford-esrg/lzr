package handshakes

import "github.com/stanford-esrg/lzr/handshakes/mongodb"

func init() {
    mongodb.RegisterHandshake()
}

