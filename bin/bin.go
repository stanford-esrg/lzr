package bin

import (
	"time"
	//"context"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sync"

	"github.com/stanford-esrg/lzr"
)

func LZRMain() {
	// create a context that can be cancelled
	//ctx, cancel := context.WithCancel(context.Background())

	start := time.Now()

	//read in config
	options, ok := lzr.Parse()
	if !ok {
		fmt.Fprintln(os.Stderr, "Failed to parse command line options, exiting.")
		return
	}

	//For CPUProfiling
	if options.CPUProfile != "" {
		f, err := os.Create(options.CPUProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	//initalize
	ipMeta := lzr.ConstructPacketStateMap(options)
	f := lzr.InitFile(options.Filename)
	lzr.InitParams()

	writingQueue := lzr.ConstructWritingQueue(options.Workers)
	pcapIncoming := lzr.ConstructPcapRoutine(options.Workers, options.IPv6Enabled)
	timeoutQueue := lzr.ConstructTimeoutQueue(options.Workers)
	retransmitQueue := lzr.ConstructRetransmitQueue(options.Workers)
	timeoutIncoming := lzr.PollTimeoutRoutine(
		&ipMeta, timeoutQueue, retransmitQueue, options.Workers, options.Timeout, options.RetransmitSec)
	incoming := lzr.ConstructIncomingRoutine(options.Workers)
	var incomingDone sync.WaitGroup
	incomingDone.Add(options.Workers)
	done := false
	writing := false

	// record to file
	go func() {
		for {
			select {
			case input := <-writingQueue:
				writing = true
				f.Record(input, options.Handshakes)
				writing = false
			}
		}
	}()
	//start all workers

	//read from zmap
	var wgForHandshakeComplete = new(sync.WaitGroup)
	wgForHandshakeComplete.Add(1)
	for i := 0; i < options.Workers; i++ {
		go func(i int) {
			//ExitCondition: incoming channel closed
			if i == options.Workers-1 {
				// a band-aid has been added to check to see if the number of items
				// in ipMeta has stayed the same for numHandshakes * timeout*2 time
				// if so, close. there is some non-deterministic
				// infinite loop condition that needs to be fixed in which ipMeta
				// does not empty all the way
				var infiniteLoop = false
				var lastIpMetaSize = ipMeta.Count()
				var intervalLoop = options.Timeout * lzr.NumHandshakes() * 2
				go func() {
					defer wgForHandshakeComplete.Done()
					for {
						startTime := time.Now()
						for time.Since(startTime) < time.Duration(intervalLoop)*time.Second {
							if ipMeta.HasUpdates() {
								startTime = time.Now()
								ipMeta.ResetUpdates()
								continue
							}
							time.Sleep(time.Duration(intervalLoop) * time.Millisecond)
						}
						ipMetaSize, isIpMetaChanged := ipMeta.CountAndHasUpdates()
						if lastIpMetaSize == ipMetaSize && !isIpMetaChanged {
							fmt.Fprintln(os.Stderr, "Infinite Loop, Breaking.")
							infiniteLoop = true
							return
						} else if ipMetaSize == 0 && isIpMetaChanged {
							fmt.Fprintln(os.Stderr, "Everything has been processed, Breaking.")
						} else {
							lastIpMetaSize = ipMeta.Count()
						}
					}
				}()
				for {
					if ipMeta.IsEmpty() /*|| infiniteLoop*/ {
						done = true
						break
					}
					if infiniteLoop {
						done = true
						break
					}
					//slow down to prevent CPU busy looping
					time.Sleep(1 * time.Second)
					fmt.Fprintln(os.Stderr, "Processing:", ipMeta.Count())
				}
			}

			for input := range incoming {
				if lzr.ReadZMap() {
					toACK := true
					toPUSH := false
					lzr.SendAck(options, input, &ipMeta, timeoutQueue,
						retransmitQueue, writingQueue, toACK, toPUSH, lzr.ACK)
				} else {
					lzr.SendSyn(input, &ipMeta, timeoutQueue)
				}
				ipMeta.FinishProcessing(input)
			}
			incomingDone.Done()
			return
		}(i)
	}

	//read from pcap
	for i := 0; i < options.Workers; i++ {
		go func(i int) {
			for input := range pcapIncoming {
				//fmt.Println("pcap incoming")
				//fmt.Println(input)
				inMap, startProcessing := ipMeta.IsStartProcessing(input)
				//if not in map, return
				if !inMap {
					continue
				}
				//if another thread is processing, put input back
				if !startProcessing {
					pcapIncoming <- input
					continue
				}
				lzr.HandlePcap(options, input, &ipMeta, timeoutQueue,
					retransmitQueue, writingQueue)
				ipMeta.FinishProcessing(input)
				//fmt.Println("finished pcap:")
				//fmt.Println(input)
			}
		}(i)
	}

	//read from timeout
	for i := 0; i < options.Workers; i++ {
		go func(i int) {

			for input := range timeoutIncoming {
				inMap, startProcessing := ipMeta.IsStartProcessing(input)
				//if another thread is processing, put input back
				//if not in map, return
				if !inMap {
					continue
				}
				if !startProcessing {
					timeoutIncoming <- input
					continue
				}
				lzr.HandleTimeout(options, input, &ipMeta, timeoutQueue, retransmitQueue, writingQueue)
				ipMeta.FinishProcessing(input)
			}
		}(i)
	}

	//exit gracefully when done
	incomingDone.Wait()

	//wait for the cooldown to make sure to receive incoming packets
	fmt.Fprintf(os.Stderr, "In cooldown for %d seconds to receive SYN ACK packets on the wire\n", options.Cooldown)
	time.Sleep(time.Duration(options.Cooldown) * time.Second)

	fmt.Fprintln(os.Stderr, "incoming is done, now waiting for handshakes timeout to be expired")
	wgForHandshakeComplete.Wait()
	fmt.Fprintln(os.Stderr, "handshakes are completed")

	for {
		if done && len(writingQueue) == 0 && !writing {
			if options.MemProfile != "" {
				f, err := os.Create(options.MemProfile)
				if err != nil {
					log.Fatal(err)
				}
				pprof.WriteHeapProfile(f)
				f.Close()
			}
			//closing file
			f.F.Flush()
			t := time.Now()
			elapsed := t.Sub(start)
			lzr.Summarize(elapsed)
			return
		}
	}

} //end of main
