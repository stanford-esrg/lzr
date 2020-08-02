# LZR
LZR is a shim that works in tangent with [ZMap](https://github.com/zmap/zmap) to efficiently detect and fingerprint unexpected services running on unexpected ports. LZR can detect up 18 unique protocols simoultaneously with just two extra packets and can fingerprint over 35 different protocols. LZR can also be used as a filter between ZMap and [ZGrab](https://github.com/zmap/zgrab2) to specify to ZGrab which L7 handshake to complete. 

To learn more about LZR read [this](todo).

## Building

###HyperACKtive
Note: If a host responds both on the expected port and on the random ephemeral port, whichever response comes first will dictate whether the host is marked as HyperACKtive. The expected port is contacted first, so unless there is some congestion which causes the packets to be delivered out of order, then the expected port is expected to answer first.  
