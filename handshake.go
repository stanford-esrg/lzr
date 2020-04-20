package lzr

var (

	handshakes map[string]Handshake
	fingerprintMap  map[string]int
)
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
	for hname, hand := range handshakes {
		tfingerprint = hand.Verify( data )
		if tfingerprint != "" {
			fingerprintMap[hname] += 1
			if fingerprint == "" {
				fingerprint += tfingerprint
			} else {
				fingerprint += ("-" + tfingerprint)
			}
		}
	}
	return fingerprint
}

func GetFingerprints() map[string]int {
	return fingerprintMap
}

func init() {
	handshakes = make( map[string]Handshake )
	fingerprintMap = make( map[string]int )
}
