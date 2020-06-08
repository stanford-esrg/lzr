# LZR
LZR: Accurately and Efficiently Identifying Service Liveness and Protocols in Internet-Wide Scans

###HyperACKtive
Note: If a host responds both on the expected port and on the random ephemeral port, whichever response comes first will dictate whether the host is marked as HyperACKtive. The expected port is contacted first, so unless there is some congestion which causes the packets to be delivered out of order, then the expected port is expected to answer first.  
