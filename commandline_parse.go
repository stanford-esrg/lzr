package main
import (
  "flag"
  "fmt"
  "time"
)

var (

    filename    *string
    workers     *int
)

type options struct {

    Filename   string
    Workers    int
}


// Basic flag declarations are available for string, integer, and boolean options.
func init() {
  //port = flag.Int("port", 3000, "port number")
  fname := "default_"+string(time.Now().Format("20060102150405"))+".json"
  filename = flag.String("f", fname , "json file name")
  workers = flag.Int("w", 1000 , "number of worker threads for each channel")
}

func parse() *options {

    flag.Parse()
    opt := &options{
        Filename: *filename,
        Workers: *workers,
    }
    fmt.Println("Writing results to file: ", *filename)
    fmt.Println("Number of worker threads for each channel: ", *workers)
    //fmt.Println("port:", *port)
    return opt
}
