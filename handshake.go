package lzr

import (
	"fmt"
	"strings"
)

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

// implement a hiearchy where when responses match
// for two fingerprints, we choose the more specific one
// e.g., protocols implemented on http
func hiearchizeFingerprint( fingerprint string ) string {

	if strings.Contains( fingerprint, "ipp" ) {
		return "ipp"
	}else if (strings.Contains( fingerprint, "dns") &&
		strings.Contains( fingerprint, "http")) {
		return "http"
	}else if (strings.Contains( fingerprint, "ssh") &&
		strings.Contains( fingerprint, "http")) {
		return "http"
	}else if (strings.Contains( fingerprint, "ftp") &&
		strings.Contains( fingerprint, "http")) {
		return "http"
	}else if (strings.Contains( fingerprint, "ftp") &&
		strings.Contains( fingerprint, "ssh")) {
		return "ssh"
	}else {
		fmt.Println("WARNING: NEW MULTI-FINGERPRINT:", fingerprint)
		return fingerprint
	}

}

//TODO: implement some type of hiearchy for labeling
func fingerprintResponse( data string ) string {
	fingerprint := ""
	tfingerprint := ""
	multiprint := false
	for hname, hand := range handshakes {
		tfingerprint = hand.Verify( data )
		if tfingerprint != "" {
			fingerprintMap[hname] += 1
			if fingerprint == "" {
				fingerprint += tfingerprint
			} else {
				multiprint = true
				fingerprint += ("-" + tfingerprint)
			}
		}
	}
	if multiprint {
		return hiearchizeFingerprint( fingerprint )
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
