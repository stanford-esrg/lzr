package handshakes

import "github.com/stanford-esrg/lzr/handshakes/mysql"

func init() {
	mysql.RegisterHandshake()
}

