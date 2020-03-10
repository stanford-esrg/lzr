package main
import (
  "flag"
  "fmt"
  "time"
)

var (

    filename    *string
    workers     *int
    timeout     *int
)

type options struct {

    Filename   string
    Workers    int
    Timeout    int
}


// Basic flag declarations are available for string, integer, and boolean options.
func init() {
  //port = flag.Int("port", 3000, "port number")
  fname := "default_"+string(time.Now().Format("20060102150405"))+".json"
  filename = flag.String("f", fname , "json file name")
  workers = flag.Int("w", 1000 , "number of worker threads for each channel")
  timeout = flag.Int("t", 1 , "number of seconds to wait in timeout queue")
}

func parse() *options {

    flag.Parse()
    opt := &options{
        Filename: *filename,
        Workers: *workers,
        Timeout: *timeout,
    }
    fmt.Println("Writing results to file: ", *filename)
    fmt.Println("Worker threads: ", *workers)
    fmt.Println("Timeout: ", *timeout)
    //fmt.Println("port:", *port)
    return opt
}
