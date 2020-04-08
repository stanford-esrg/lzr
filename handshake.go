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
	h, ok := handshakes[name]
	if !ok {
		panic(name + " not found")
	}
	return h
}

//TODO: implement some type of hiearchy for labeling
func fingerprintResponse( data string ) string {
	fingerprint := ""
	tfingerprint := ""
	for _, hand := range handshakes {
		tfingerprint = hand.Verify( data )
		if tfingerprint != "" {
			if fingerprint == "" {
				fingerprint += tfingerprint
			} else {
				fingerprint += ("-" + tfingerprint)
			}
		}
	}
	return fingerprint
}


func init() {
	handshakes = make( map[string]Handshake )
}
