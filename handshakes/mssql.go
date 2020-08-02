package handshakes

import "github.com/stanford-esrg/lzr/handshakes/mssql"

func init() {
	mssql.RegisterHandshake()
}

