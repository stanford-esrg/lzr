package main

import (
    "time"
    "context"
    "golang.org/x/sync/semaphore"
    //"fmt"
)


//TODO: move vars to appropriate places
var (
    device       string = "ens8"
    snapshot_len int32  = 1024
    promiscuous  bool   = false
    err          error
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
    timeoutIncoming := pollTimeoutRoutine( &ipMeta,timeoutQueue, options.Workers, options.Timeout )
    timeoutEmpty := true

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


	//read from both zmap and pcap
	for {
		select {
			case input, ok := <-zmapIncoming:
                //ExitCondition:
                //done reading from zmap (channel closed)
                //TODO: timeoutQ is empty ??????
                //eventually check: all locks are released (no more jobs runnign)
                if !ok && timeoutEmpty {
                    time.Sleep( 1 * time.Second )
                    return
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
			case input := <-pcapIncoming:

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
            case input, _ := <-timeoutIncoming:
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
                    handleTimeout( input, &ipMeta, &timeoutQueue, &writingQueue, f )
                    ipMeta.finishProcessing( &input )
                }()

			default:
                continue
		}
	}
} //end of main
