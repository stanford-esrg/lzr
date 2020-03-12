package main

import (
    "net/http"
    "net/http/httputil"
    "strings"
)

func getData( dst string ) []byte {

        req, _ := http.NewRequest("GET","/",nil)
        req.Host =  dst
        req.Header.Add("Host",dst)
        req.Header.Set("User-Agent","Mozilla/5.0 zgrab/0.x")
        req.Header.Set("Accept","*/*")
        req.Header.Set("Accept-Encoding","gzip")
        data, _ := httputil.DumpRequest(req, false)
	return data
}

func fingerprintResponse( data string ) string {

    if strings.Contains( data, "HTTP" ) {
        return "HTTP"
    } else {
        return ""
    }

}
