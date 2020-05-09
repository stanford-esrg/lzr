package lzr
import (
  "flag"
  "fmt"
  "time"
  "strings"
)

var (

    filename			*string
	debug				*bool
	haf					*bool
	pushDOnly			*bool
	forceAllHandshakes	*bool
	feedZGrab			*bool
    workers				*int
    timeout				*int
    retransmitSec		*int
    retransmitNum		*int
	cpuprofile			*string
	memprofile			*string
	handshake			*string
)

type options struct {

    Filename			string
	Debug				bool
	Haf					bool
	PushDOnly			bool
	ForceAllHandshakes	bool
	FeedZGrab			bool
    Workers				int
    Timeout				int
    RetransmitSec		int
    RetransmitNum		int
	CPUProfile			string
	MemProfile			string
	Handshakes			[]string
}


// Basic flag declarations are available for string, integer, and boolean options.
func init() {
  fname := "default_"+string(time.Now().Format("20060102150405"))+".json"
  filename = flag.String("f", fname , "json file name")

  debug = flag.Bool("d", false, "debug printing on (defaut off)")
  haf = flag.Bool("haf", true, "HyperACKtive filtering off (default on)")
  pushDOnly = flag.Bool("pushDataOnly", false, "Don't attach data to ack but rather to push only (default off)")
  forceAllHandshakes = flag.Bool("forceAllHandshakes", false, "Complete all handshakes even if data is returned early on (default off). This also turns off HyperACKtive filtering.")
  feedZGrab = flag.Bool("feedZGrab", false, "send to zgrab ip and fingerprint (default off)")
  workers = flag.Int("w", 1 , "number of worker threads for each channel")
  timeout = flag.Int("t", 5, "number of seconds to wait in timeout queue for last retransmission")
  retransmitSec = flag.Int("rt", 1 , "number of seconds until re-transmitting packet")
  retransmitNum = flag.Int("rn", 1 , "number of data packets to re-transmit")
  cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
  memprofile = flag.String("memprofile", "", "write memory profile to this file")
  handshake = flag.String("handshakes", "http" , "handshakes to scan with")
}

func Parse() *options {

    flag.Parse()
    opt := &options{
        Filename: *filename,
		Debug: *debug,
		Haf: *haf,
		FeedZGrab: *feedZGrab,
		PushDOnly: *pushDOnly,
		ForceAllHandshakes: *forceAllHandshakes,
        Workers: *workers,
        Timeout: *timeout,
        RetransmitSec: *retransmitSec,
        RetransmitNum: *retransmitNum,
		CPUProfile: *cpuprofile,
		MemProfile: *memprofile,
		Handshakes: make([]string, strings.Count(*handshake,",")+1),
    }
	if !strings.Contains( *handshake, ",")	{
		opt.Handshakes[0] = *handshake
	} else {
		i := 0
		for _, h := range strings.Split( *handshake, "," ) {
			opt.Handshakes[i] = h
			i += 1
		}
	}

	if *forceAllHandshakes {
		*haf = false
	}

    fmt.Println("++Writing results to file:", *filename)
    fmt.Println("++Handshakes:", *handshake)
	if *memprofile != "" {
		fmt.Println("++Writing memprofile to file:", *memprofile)
	}
	if *cpuprofile != "" {
		fmt.Println("++Writing cpuprofile to file:", *cpuprofile)
	}
	if *debug {
		fmt.Println("++Debug turned on")
	}
	if *haf {
		fmt.Println("++HyperACKtiveFiltering turned on")
	}
	if *feedZGrab {
		fmt.Println("++Feeding ZGrab with fingerprints")
	}
	if *pushDOnly {
		fmt.Println("++Sending Data only with Push Flag (not in ack)")
	}
	if *forceAllHandshakes {
		fmt.Println("++Force completing all handshakes")
	}
    fmt.Println("++Worker threads:", *workers)
    fmt.Println("++Timeout Interval (s):", *timeout)
    fmt.Println("++Retransmit Interval (s):", *retransmitSec)
    fmt.Println("++Number of Retransmitions:", *retransmitNum)
    //fmt.Println("port:", *port)
    return opt
}

func DebugOn() bool {
	return *debug
}

func FeedZGrab() bool {
	return *feedZGrab
}

func HyperACKtiveFiltering() bool {
	return *haf
}

func PushDOnly() bool {
	return *pushDOnly
}

func ForceAllHandshakes() bool {
	return *forceAllHandshakes
}
