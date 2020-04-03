package main

import (
	"lzr/bin"
	_ "lzr/handshakes"
)

// main wraps the "true" main, bin.LZRMain()
// after importing all handshake modules
func main() {
	bin.LZRMain()
}
