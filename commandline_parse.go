package main
import (
  "flag"
  "fmt"
)

var (
  port *int
)

// Basic flag declarations are available for string, integer, and boolean options.
func init() {
  port = flag.Int("port", 3000, "an int")
}

func parse() int {

    flag.Parse()
    fmt.Println("port:", *port)
    return *port
}
