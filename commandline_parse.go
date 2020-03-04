package main
import (
  "flag"
  "fmt"
  "time"
)

var (

    filename    *string
)

type options struct {

    Filename   string

}


// Basic flag declarations are available for string, integer, and boolean options.
func init() {
  //port = flag.Int("port", 3000, "port number")
  fname := "default_"+string(time.Now().Format("20060102150405"))+".json"
  filename = flag.String("f", fname , "json file name")
}

func parse() *options {

    flag.Parse()
    opt := &options{
        Filename: *filename,
    }
    fmt.Println("writing results to file ", *filename)
    //fmt.Println("port:", *port)
    return opt
}
