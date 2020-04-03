package lzr

var handshakes map[string]Handshake

type Handshake interface {

    //get Data to send in first packet
    GetData( dst string ) []byte

}

func AddHandshake( name string, h Handshake ) {

	handshakes[ name ] = h

}

func GetHandshake( name string ) Handshake {
	return handshakes[name]
}

func init() {
	handshakes = make( map[string]Handshake )
}
