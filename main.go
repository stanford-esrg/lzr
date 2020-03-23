package main

import (
    //"time"
    "context"
    "golang.org/x/sync/semaphore"
    "fmt"
)



func main() {

    ctx := context.TODO()

    //read in config 
    options := parse()

	//initalize
	ipMeta := constructPacketStateMap()
    f := initFile( options.Filename )
    sem := semaphore.NewWeighted( int64(options.Workers) )

    writingQueue := constructWritingQueue( options.Workers )
    zmapIncoming := constructZMapRoutine( options.Workers )
    pcapIncoming := constructPcapRoutine( options.Workers )
    timeoutQueue := constructTimeoutQueue( options.Workers )
    timeoutIncoming,timeoutEmptyChannel := pollTimeoutRoutine(
        &ipMeta,timeoutQueue, options.Workers, options.Timeout )
    done := false

    // record to file
    go func() {
        for {
            select {
                case input := <-writingQueue:
                    f.record( input )
                default:
                    continue
                }
        }
    }()
    //read from zmap
    go func() {
	    for {
		    select {
			case input, ok := <-zmapIncoming:
                //ExitCondition: zmap channel closed
                if !ok {
                    //TODO: change to see if metaMap is empty
                    if *timeoutEmptyChannel {
                        done=true
                    }
                    continue
                }
                if err := sem.Acquire(ctx, 1); err != nil {
                    continue
                }
                go func() {
                    defer sem.Release(1)
                    // not checking if another thread is processing since we are 
                    // assuming that repeats are being filtered at zmap 
                    // and thus IPs cannot be in ipMeta b4 zmap adds to it
				    ackZMap( input, &ipMeta, &timeoutQueue, &writingQueue, f )
                }()
            default:
                continue
            }
        }
    }()
    //read from pcap
    go func() {
        for {
            select {
			case input := <-pcapIncoming:
                if input.Saddr == "104.16.131.21" {
                    fmt.Println(input)
                }
                //fmt.Println( packet )
                if err := sem.Acquire(ctx, 1); err != nil {
                    continue
                }
                go func() {
                    defer sem.Release(1)
                    inMap, processing := ipMeta.isProcessing( &input )
                    //if another thread is processing, put input back
                    if processing {
                        pcapIncoming <- input
                        return
                    }
                    //if not in map, return
                    if !inMap {
                        return
                    }
                    ipMeta.startProcessing( &input )
				    handlePcap( input, &ipMeta, &timeoutQueue, &writingQueue, f )
                    ipMeta.finishProcessing( &input )
                }()
            default:
                continue
            }
        }
    }()
    //read from timeout
    go func() {

        for {
            select {
            case input, _ := <-timeoutIncoming:
                if err := sem.Acquire(ctx, 1); err != nil {
                    continue
                }
                go func() {
                    defer sem.Release(1)
                    inMap, processing := ipMeta.isProcessing( &input )
                    //if another thread is processing, put input back
                    if processing {
                        timeoutIncoming <- input
                        return
                    }
                    //if not in map, return
                    if !inMap {
                        return
                    }
                    ipMeta.startProcessing( &input )
                    handleTimeout( input, &ipMeta, &timeoutQueue, &writingQueue, f )
                    ipMeta.finishProcessing( &input )
                }()

			default:
                continue
		    }
        }
    }()

    //exit gracefull when done
    for {
        if done {
            return
        }
    }



} //end of main
