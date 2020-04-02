package main

import (
    "time"
    //"context"
    //"golang.org/x/sync/semaphore"
	"runtime/pprof"
	"sync"
	"os"
	"log"
)

var (
    PARTITIONS = 2
)

func main() {
    // create a context that can be cancelled
    //ctx, cancel := context.WithCancel(context.Background())

    //read in config 
    options := parse()

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
	ipMeta := constructPacketStateMap()
    f := initFile( options.Filename )
    //sem := semaphore.NewWeighted( int64(options.Workers) )

    writingQueue := constructWritingQueue( options.Workers )
    zmapIncoming := constructZMapRoutine( options.Workers )
    pcapIncoming, pcapQueue := constructPcapRoutine( options.Workers ) //pcapQueue
    timeoutQueue := constructTimeoutQueue( options.Workers )
    timeoutIncoming := pollTimeoutRoutine(
        &ipMeta,timeoutQueue, options.Workers, options.Timeout )
    done := false

    // record to file
    go func() {
        for {
            select {
                case input := <-writingQueue:
					//fmt.Println("writing")
                    f.record( input )
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
	zmapDone.Add(options.Workers/PARTITIONS)
    for i := 0; i < options.Workers/PARTITIONS; i ++ {
        go func( i int ) {
	        for input := range zmapIncoming {
				        ackZMap( input, &ipMeta, &timeoutQueue, &writingQueue, f )
                        ipMeta.finishProcessing( &input )
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
    for i := 0; i < options.Workers/PARTITIONS; i ++ {
        go func() {
            for input := range pcapIncoming {
                        inMap, processing := ipMeta.isStartProcessing( &input )
                        //if not in map, return
                        if !inMap {
                            continue
                        }
                        //if another thread is processing, put input back
                        if !processing {
                            pcapQueue <- input
                            continue
                        }
				        handlePcap( input, &ipMeta, &timeoutQueue, &writingQueue, f )
                        ipMeta.finishProcessing( &input )
            }
        }()
    }

    //read from timeout
    go func() {

        for input := range timeoutIncoming {
                    inMap, processing := ipMeta.isStartProcessing( &input )
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
                    handleTimeout( input, &ipMeta, &timeoutQueue, &writingQueue, f )
                    ipMeta.finishProcessing( &input )
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
            return
       }
    }



} //end of main
