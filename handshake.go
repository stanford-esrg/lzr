package lzr

var handshakes map[string]Handshake

type Handshake interface {

    //get Data to send in first packet
    GetData( dst string ) []byte
	//verify the protocol from response
	Verify( data string)  string

}

func AddHandshake( name string, h Handshake ) {

	handshakes[ name ] = h

}

func GetHandshake( name string ) Handshake {
	return handshakes[name]
}

func fingerprintResponse( data string ) string {
	var fingerprint string
	for _, hand := range handshakes {
		fingerprint = hand.Verify( data )
		if fingerprint != "" {
			return fingerprint
		}
	}
	return ""
}


func init() {
	handshakes = make( map[string]Handshake )
}
