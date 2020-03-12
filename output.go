package main

import (
    "encoding/json"
    "log"
    "os"
    "sync"
)

type output_file struct {

    F    *os.File
    MLock    sync.RWMutex

}


func ( f *output_file ) record( packet packet_metadata ) {

    out, _ := json.Marshal( packet )
    f.MLock.Lock()
    _,err = (f.F).WriteString( string(out) )
    _,err = (f.F).WriteString( "\n" )
    f.MLock.Unlock()
    if err != nil {
        f.F.Close()
		log.Fatal(err)
	}
    return
}


func initFile( fname string ) *output_file {
    f, err := os.OpenFile( fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644 )

    if err != nil {
		log.Fatal(err)
    }

    o := &output_file{
        F: f,
        MLock: sync.RWMutex{},
    }

    return o
}

//TODO: need to figure out when to close
