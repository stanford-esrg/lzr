package main

import (
    "time"
    //"context"
    //"golang.org/x/sync/semaphore"
	"runtime/pprof"
	"os"
	"log"
    "fmt"
)

var (
    PARTITIONS = 2
)

func main() {
    // create a context that can be cancelled
    //ctx, cancel := context.WithCancel(context.Background())

    //read in config 
    options := parse()

	time.Sleep(2*time.Second)
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
    pcapIncoming, pcapQueue := constructPcapRoutine( options.Workers )
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
    //start all workers
    //read from zmap
    for i := 0; i < options.Workers/PARTITIONS; i ++ {
        go func( i int ) {
	        for {
		        select {
			        case input, ok := <-zmapIncoming:
                        //ExitCondition: zmap channel closed
                        if !ok {
                            if ipMeta.IsEmpty() {
                                done=true
                                return
                            }
							fmt.Println("state map: ",ipMeta.Count())
							fmt.Println("pcapInc: ", len(pcapIncoming))
							fmt.Println("pcapQ: ", len(pcapQueue))
							time.Sleep(1*time.Second)
                            continue
                        }

                        // not checking if another thread is processing since we are 
                        // assuming that repeats are being filtered at zmap 
                        // and thus IPs cannot be in ipMeta b4 zmap adds to it
				        ackZMap( input, &ipMeta, &timeoutQueue, &writingQueue, f )
                    case <-time.After(2 * time.Second):
                        //fmt.Println(i,": Something wrong with reading from zmap")
                        continue
                }
            }
        }(i)
    }
    //read from pcap
    for i := 0; i < options.Workers/PARTITIONS; i ++ {
        go func() {
            for {
                select {
                   case input := <-pcapIncoming:
                        if input.Saddr == "104.16.128.199" {
                            fmt.Println(input)
                        }
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
						fmt.Println("pcap processing...",input.Saddr)
				        handlePcap( input, &ipMeta, &timeoutQueue, &writingQueue, f )
                        ipMeta.finishProcessing( &input )
                    case <-time.After(2 * time.Second):
                        //fmt.Println("Something wrong with reading from pcap")
                        continue
                }
            }
        }()
    }

    //read from timeout
    go func() {

        for {
            select {
            case input, _ := <-timeoutIncoming:
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
					fmt.Println("tout processing...",input.Saddr)
                    handleTimeout( input, &ipMeta, &timeoutQueue, &writingQueue, f )
                    ipMeta.finishProcessing( &input )

                case <-time.After(2 * time.Second):
                    //fmt.Println("Something wrong with reading from timeout")
                    continue
		    }
        }
    }()

    //exit gracefully when done
	//OR for debugging, within 5 seconds
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
       if done {
			fmt.Println(len(pcapIncoming))
            return
       }
    }



} //end of main
