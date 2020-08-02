package handshakes

import "github.com/stanford-esrg/lzr/handshakes/kubernetes"

func init() {
	kubernetes.RegisterHandshake()
}

