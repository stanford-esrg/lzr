package handshakes

import "lzr/handshakes/kubernetes"

func init() {
	kubernetes.RegisterHandshake()
}

