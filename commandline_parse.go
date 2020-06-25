package lzr
import (
  "flag"
  "fmt"
  "os"
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
	handshakeArr		[]string
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

  debug = flag.Bool("d", false, "debug printing on")
  haf = flag.Bool("haf", true, "HyperACKtive filtering off")
  pushDOnly = flag.Bool("pushDataOnly", false, "Don't attach data to ack but rather to push only")
  forceAllHandshakes = flag.Bool("forceAllHandshakes", false, "Complete all handshakes even if data is returned early on. This also turns off HyperACKtive filtering.")
  feedZGrab = flag.Bool("feedZGrab", false, "send to zgrab ip and fingerprint")
  workers = flag.Int("w", 1 , "number of worker threads for each channel")
  timeout = flag.Int("t", 5, "number of seconds to wait in timeout queue for last retransmission")
  retransmitSec = flag.Int("rt", 1 , "number of seconds until re-transmitting packet")
  retransmitNum = flag.Int("rn", 1 , "number of data packets to re-transmit")
  cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
  memprofile = flag.String("memprofile", "", "write memory profile to this file")
  handshake = flag.String("handshakes", "http" , "handshakes to scan with")
}



func Parse() (*options,bool) {

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

		_, ok := GetHandshake(*handshake)

		if !ok {
			fmt.Fprintln(os.Stderr,"--Handshake not found:", *handshake)
			return nil,false
		}

		opt.Handshakes[0] = *handshake
	} else {
		i := 0
		for _, h := range strings.Split( *handshake, "," ) {

			_, ok := GetHandshake(h)
			if !ok {
				fmt.Fprintln(os.Stderr,"--Handshake not found:", h)
				return nil,false
			}

			opt.Handshakes[i] = h
			i += 1
		}
	}
	handshakeArr = opt.Handshakes

	if *forceAllHandshakes {
		*haf = false
	}

    fmt.Fprintln(os.Stderr,"++Writing results to file:", *filename)
    fmt.Fprintln(os.Stderr,"++Handshakes:", *handshake)
	if *memprofile != "" {
		fmt.Fprintln(os.Stderr,"++Writing memprofile to file:", *memprofile)
	}
	if *cpuprofile != "" {
		fmt.Fprintln(os.Stderr,"++Writing cpuprofile to file:", *cpuprofile)
	}
	if *debug {
		fmt.Fprintln(os.Stderr,"++Debug turned on")
	}
	if *haf {
		fmt.Fprintln(os.Stderr,"++HyperACKtiveFiltering turned on")
	}
	if *feedZGrab {
		fmt.Fprintln(os.Stderr,"++Feeding ZGrab with fingerprints")
	}
	if *pushDOnly {
		fmt.Fprintln(os.Stderr,"++Sending Data only with Push Flag (not in ack)")
	}
	if *forceAllHandshakes {
		fmt.Fprintln(os.Stderr,"++Force completing all handshakes")
	}
    fmt.Fprintln(os.Stderr,"++Worker threads:", *workers)
    fmt.Fprintln(os.Stderr,"++Timeout Interval (s):", *timeout)
    fmt.Fprintln(os.Stderr,"++Retransmit Interval (s):", *retransmitSec)
    fmt.Fprintln(os.Stderr,"++Number of Retransmitions:", *retransmitNum)
    //fmt.Fprintln(os.Stderr,"port:", *port)
    return opt,true
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

func GetAllHandshakes()  []string {
	return handshakeArr
}
