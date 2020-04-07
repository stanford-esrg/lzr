package kubernetes

import (
    "net/http"
    "net/http/httputil"
	"strings"
	//"fmt"
	"lzr"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {

        req, _ := http.NewRequest("GET","/api",nil)
        req.Host =  dst
        req.Header.Add("Host",dst)
        req.Header.Set("User-Agent","Mozilla/5.0 zgrab/0.x")
        req.Header.Set("Accept","*/*")
        //req.Header.Set("Content-Type","application/ipp")
        req.Header.Set("Accept-Encoding","gzip")
        data, _ := httputil.DumpRequest(req, false)
    return data
}

func (h *HandshakeMod) Verify( data string ) string {

	if	 strings.Contains( data, "200 OK" ) &&
		 strings.Contains( data, "APIVersions" ){
         return "kubernetes"
	}
	return ""

}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "kubernetes", &h )
}

