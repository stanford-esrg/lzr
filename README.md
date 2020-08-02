LZR
=========

LZR is a shim that works in tangent with [ZMap](https://github.com/zmap/zmap) to efficiently detect and fingerprint unexpected services running on unexpected ports. LZR can detect up 18 unique protocols simoultaneously with just two extra packets and can fingerprint over 35 different protocols. LZR can also be used as a filter between ZMap and [ZGrab](https://github.com/zmap/zgrab2) to specify to ZGrab which L7 handshake to complete. 

To learn more about LZR read [this](todo).

## Building

Set up `$GOPATH` (see https://golang.org/doc/code.html)
```
$ go get github.com/stanford-esrg/LZR
$ cd $GOPATH/src/github.com/stanford-esrg/LZR
```

LZR intercepts connections which ZMap opens; in order to ensure that the kernel does not interfere with LZR, LZR requires a source-ip to be specified for which the kernel drops all RSTs for traffic targeted to the source-ip. The chosen source-ip&mdash;which both ZMap and LZR will use&mdash;should be passed in as a parameter to make, so the appropriate iptables rule can be set.
```
$ make all source-ip=256.256.256.256/32
```

## Usage

#### HyperACKtive
Note: If a host responds both on the expected port and on the random ephemeral port, whichever response comes first will dictate whether the host is marked as HyperACKtive. The expected port is contacted first, so unless there is some congestion which causes the packets to be delivered out of order, then the expected port is expected to answer first.  
