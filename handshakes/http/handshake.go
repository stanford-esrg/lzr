package http

import (
    "net/http"
    "net/http/httputil"
	"lzr"
)

// Handshake implements the lzr.Handshake interface
type Handshake struct {
}

func (h *Handshake) getData( dst string ) []byte {

        req, _ := http.NewRequest("GET","/",nil)
        req.Host =  dst
        req.Header.Add("Host",dst)
        req.Header.Set("User-Agent","Mozilla/5.0 zgrab/0.x")
        req.Header.Set("Accept","*/*")
        req.Header.Set("Accept-Encoding","gzip")
        data, _ := httputil.DumpRequest(req, false)
    return data
}

func RegisterHandshake() {
	var h lzr.Handshake
	lzr.AddHandshake( "http",h )
}

