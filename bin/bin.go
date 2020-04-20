package bin

import (
    "time"
    //"context"
	"runtime/pprof"
	"sync"
	"os"
	"log"
	"lzr"
	//"fmt"
)

func LZRMain() {
    // create a context that can be cancelled
    //ctx, cancel := context.WithCancel(context.Background())

    //read in config 
    options := lzr.Parse()

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
	ipMeta := lzr.ConstructPacketStateMap()
    f := lzr.InitFile( options.Filename )

    writingQueue := lzr.ConstructWritingQueue( options.Workers )
    zmapIncoming := lzr.ConstructZMapRoutine( options.Workers )
    pcapIncoming, pcapQueue := lzr.ConstructPcapRoutine( options.Workers )
    timeoutQueue := lzr.ConstructTimeoutQueue( options.Workers )
    timeoutIncoming := lzr.PollTimeoutRoutine(
        &ipMeta,timeoutQueue, options.Workers, options.Timeout )
    done := false
	writing := false


    // record to file
    go func() {
        for {
            select {
                case input := <-writingQueue:
					writing = true
                    f.Record( *input, options.Handshakes )
					writing = false
                }
        }
    }()
    //start all workers
    //read from zmap
	var zmapDone sync.WaitGroup
	zmapDone.Add(options.Workers)
    for i := 0; i < options.Workers; i ++ {
        go func( i int ) {
	        for input := range zmapIncoming {
				        lzr.SendAck( options.Handshakes, input, &ipMeta, &timeoutQueue, &writingQueue )
                        ipMeta.FinishProcessing( input )
            }

            //ExitCondition: zmap channel closed
			for {
				if ipMeta.IsEmpty() {
					done=true
					zmapDone.Done()
					return
				}
				time.Sleep(1*time.Second)
            }
        }(i)
    }
    //read from pcap
    for i := 0; i < options.Workers; i ++ {
        go func() {
            for input := range pcapIncoming {
                        inMap, processing := ipMeta.IsStartProcessing( input )
                        //if not in map, return
                        if !inMap {
                            continue
                        }
                        //if another thread is processing, put input back
                        if !processing {
                            pcapQueue <- input
                            continue
                        }
				        lzr.HandlePcap(options.Handshakes, input, &ipMeta, &timeoutQueue, &writingQueue )
                        ipMeta.FinishProcessing( input )
            }
        }()
    }

    //read from timeout
    go func() {

        for input := range timeoutIncoming {
                    inMap, processing := ipMeta.IsStartProcessing( input )
                    //if another thread is processing, put input back
                    //if not in map, return
                    if !inMap {
                        continue
                    }
                    if !processing {
                        timeoutQueue <- input
                        continue
                    }
                    lzr.HandleTimeout( options.Handshakes, input, &ipMeta, &timeoutQueue, &writingQueue )
                    ipMeta.FinishProcessing( input )
		    }
    }()

    //exit gracefully when done
	//OR for debugging, within 5 seconds
	zmapDone.Wait()
    for {
       if done && len(writingQueue) == 0 && !writing {
			//TODO: need to properly close file
				if options.MemProfile != "" {
					f, err := os.Create(options.MemProfile)
					if err != nil {
						log.Fatal(err)
					}
					pprof.WriteHeapProfile(f)
					f.Close()
				}

            return
       }
    }



} //end of main
