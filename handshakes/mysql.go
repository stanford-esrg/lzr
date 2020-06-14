package handshakes

import "lzr/handshakes/mysql"

func init() {
	mysql.RegisterHandshake()
}

