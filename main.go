package main

import (
    "time"
    //"context"
    //"golang.org/x/sync/semaphore"
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

	//initalize
	ipMeta := constructPacketStateMap()
    f := initFile( options.Filename )
    //sem := semaphore.NewWeighted( int64(options.Workers) )

    writingQueue := constructWritingQueue( options.Workers )
    zmapIncoming := constructZMapRoutine( options.Workers )
    pcapIncoming := constructPcapRoutine( options.Workers )
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
                        inMap, processing := ipMeta.isProcessing( &input )
                        //if another thread is processing, put input back
                        if processing {
                            pcapIncoming <- input
                            continue
                        }
                        //if not in map, return
                        if !inMap {
                            continue
                        }
                        ipMeta.startProcessing( &input )
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
                    inMap, processing := ipMeta.isProcessing( &input )
                    //if another thread is processing, put input back
                    if processing {
                        timeoutQueue <- input // Incoming or Q to avoid dlock??
                        continue
                        //return
                    }
                    //if not in map, return
                    if !inMap {
                        continue
                        //return
                    }
                    ipMeta.startProcessing( &input )
                    handleTimeout( input, &ipMeta, &timeoutQueue, &writingQueue, f )
                    ipMeta.finishProcessing( &input )

                case <-time.After(2 * time.Second):
                    //fmt.Println("Something wrong with reading from timeout")
                    continue
		    }
        }
    }()

    //exit gracefully when done
    for {
        if done {
            return
        }
    }



} //end of main
