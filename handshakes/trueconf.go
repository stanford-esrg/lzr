package handshakes

import (
	"github.com/stanford-esrg/lzr"
	"github.com/stanford-esrg/lzr/handshakes/trueconf"
)

func init() {
	lzr.AddHandshake("trueconf", trueconf.Handshake{})
}
