package http

import (
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/stanford-esrg/lzr"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData(dst string) []byte {

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = dst
	if lzr.IPv6Enabled() {
		req.Header.Add("Host", "["+dst+"]")
	} else {
		req.Header.Add("Host", dst)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.3")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip")
	data, _ := httputil.DumpRequest(req, false)
	return data
}

func (h *HandshakeMod) Verify(data string) string {

	if !strings.Contains(data, "HTTPS") &&
		(strings.Contains(data, "HTTP") || strings.Contains(data, "html") ||
			strings.Contains(data, "HTML") || strings.Contains(data, "<h1>")) {
		return "http"
	}
	return ""

}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake("http", &h)
}
