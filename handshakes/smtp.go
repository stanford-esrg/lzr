package handshakes

import "github.com/stanford-esrg/lzr/handshakes/smtp"

func init() {
	smtp.RegisterHandshake()
}

