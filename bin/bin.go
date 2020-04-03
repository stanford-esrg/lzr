package bin

import (
    "time"
    //"context"
    //"golang.org/x/sync/semaphore"
	"runtime/pprof"
	"sync"
	"os"
	"log"
	"lzr"
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
	handshake := lzr.GetHandshake(options.Handshake) // ??? 
	ipMeta := lzr.ConstructPacketStateMap()
    f := lzr.InitFile( options.Filename )
    //sem := semaphore.NewWeighted( int64(options.Workers) )

    writingQueue := lzr.ConstructWritingQueue( options.Workers )
    zmapIncoming := lzr.ConstructZMapRoutine( options.Workers )
    pcapIncoming, pcapQueue := lzr.ConstructPcapRoutine( options.Workers ) //pcapQueue
    timeoutQueue := lzr.ConstructTimeoutQueue( options.Workers )
    timeoutIncoming := lzr.PollTimeoutRoutine(
        &ipMeta,timeoutQueue, options.Workers, options.Timeout )
    done := false

    // record to file
    go func() {
        for {
            select {
                case input := <-writingQueue:
					//fmt.Println("writing")
                    f.Record( input )
                 case <-time.After(2 * time.Second):
					//fmt.Println(": Something wrong with reading from writing")
                    continue
                }
        }
    }()
	time.Sleep(1*time.Second)
    //start all workers
    //read from zmap
	var zmapDone sync.WaitGroup
	zmapDone.Add(options.Workers)
    for i := 0; i < options.Workers; i ++ {
        go func( i int ) {
	        for input := range zmapIncoming {
				        lzr.AckZMap( handshake, input, &ipMeta, &timeoutQueue, &writingQueue, f )
                        ipMeta.FinishProcessing( &input )
            }

            //ExitCondition: zmap channel closed
			for {
				if ipMeta.IsEmpty() {
					done=true
					zmapDone.Done()
					return
				}
            }
        }(i)
    }
	time.Sleep(1*time.Second)
    //read from pcap
    for i := 0; i < options.Workers; i ++ {
        go func() {
            for input := range pcapIncoming {
                        inMap, processing := ipMeta.IsStartProcessing( &input )
                        //if not in map, return
                        if !inMap {
                            continue
                        }
                        //if another thread is processing, put input back
                        if !processing {
                            pcapQueue <- input
                            continue
                        }
				        lzr.HandlePcap(handshake, input, &ipMeta, &timeoutQueue, &writingQueue, f )
                        ipMeta.FinishProcessing( &input )
            }
        }()
    }

    //read from timeout
    go func() {

        for input := range timeoutIncoming {
                    inMap, processing := ipMeta.IsStartProcessing( &input )
                    //if another thread is processing, put input back
                    //if not in map, return
                    if !inMap {
                        continue
                        //return
                    }
                    if !processing {
                        timeoutQueue <- input 
                        continue
                        //return
                    }
                    lzr.HandleTimeout( handshake, input, &ipMeta, &timeoutQueue, &writingQueue, f )
                    ipMeta.FinishProcessing( &input )
		    }
    }()

    //exit gracefully when done
	//OR for debugging, within 5 seconds
	zmapDone.Wait()
    for {
		/*select {
			case <-time.After(30 * time.Second):

				if options.MemProfile != "" {
					f, err := os.Create(options.MemProfile)
					if err != nil {
						log.Fatal(err)
					}
					pprof.WriteHeapProfile(f)
					f.Close()
				}
				return
		}*/
       if done && len(writingQueue) == 0 {
			//TODO: need to properly close file
            return
       }
    }



} //end of main
