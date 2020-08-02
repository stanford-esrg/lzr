package handshakes

import "github.com/stanford-esrg/lzr/handshakes/ssh"

func init() {
	ssh.RegisterHandshake()
}

