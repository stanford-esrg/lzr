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

func GetHandshake( name string ) (Handshake,bool) {
	h, ok := handshakes[name]
	return h,ok
}

// implement a hiearchy where when responses match
// for two fingerprints, we choose the more specific one
// e.g., protocols implemented on http
func hiearchizeFingerprint( fingerprint string ) string {


	// prioritize for the handshake being sent
	// or for handshake which was asked to be prioritized
	// so if scanning for http ipp will return as http
	// but if scanning for ipp then http+ipp will return ipp

	req_handshakes := GetAllHandshakes()
	for  _, h := range req_handshakes {
		if strings.Contains( fingerprint, h ) {
			return h
		}
	}


	if strings.Contains( fingerprint, "ipp" ) {
		return "ipp"
	} else if strings.Contains( fingerprint, "kubernetes") {
		return "kubernetes"
	} else if (strings.Contains( fingerprint, "dns") &&
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
	//probs tls with HTTPS text
	} else if (strings.Contains( fingerprint, "tls") &&
		strings.Contains( fingerprint, "http")) {
		return "tls"
	} else {
		fmt.Println("WARNING: NEW MULTI-FINGERPRINT:", fingerprint)
		return fingerprint
	}

}

//TODO: implement some type of hiearchy for labeling
func fingerprintResponse( data string ) string {
	fingerprint := ""
	tfingerprint := ""
	multiprint := false
	for _, hand := range handshakes {
		tfingerprint = hand.Verify( data )
		if tfingerprint != "" {
			//concat fingerprints together 
			if fingerprint == "" {
				fingerprint += tfingerprint
			} else {
				multiprint = true
				fingerprint += ("-" + tfingerprint)
			}

		}
	}
	if multiprint {
		fingerprint = hiearchizeFingerprint( fingerprint )
	}
	fingerprintMap[fingerprint] += 1
	return fingerprint
}

func GetFingerprints() map[string]int {
	return fingerprintMap
}

func init() {
	handshakes = make( map[string]Handshake )
	fingerprintMap = make( map[string]int )
}
