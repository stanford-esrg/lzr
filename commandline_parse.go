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
	cpuprofile	*string
	memprofile	*string
)

type options struct {

    Filename   string
    Workers    int
    Timeout    int
	CPUProfile string
	MemProfile string
}


// Basic flag declarations are available for string, integer, and boolean options.
func init() {
  //port = flag.Int("port", 3000, "port number")
  fname := "default_"+string(time.Now().Format("20060102150405"))+".json"
  filename = flag.String("f", fname , "json file name")
  workers = flag.Int("w", 1000 , "number of worker threads for each channel")
  timeout = flag.Int("t", 1 , "number of seconds to wait in timeout queue")
  cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
  memprofile = flag.String("memprofile", "", "write memory profile to this file")
}

func parse() *options {

    flag.Parse()
    opt := &options{
        Filename: *filename,
        Workers: *workers,
        Timeout: *timeout,
		CPUProfile: *cpuprofile,
		MemProfile: *memprofile,
    }
    fmt.Println("Writing results to file: ", *filename)
	if *memprofile != "" {
		fmt.Println("Writing memprofile to file: ", *cpuprofile)
	}
	if *cpuprofile != "" {
		fmt.Println("Writing cpuprofile to file: ", *cpuprofile)
	}
    fmt.Println("Worker threads: ", *workers)
    fmt.Println("Timeout: ", *timeout)
    //fmt.Println("port:", *port)
    return opt
}
