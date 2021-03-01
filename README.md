LZR
=========

LZR quickly detects and fingerprints unexpected services running on unexpected ports by working with [ZMap](https://github.com/zmap/zmap). LZR can detect up 18 unique protocols simultaneously with just two extra packets and can fingerprint over 35 different protocols. 

To learn more about LZR's system and performance, check out the original [paper](https://lizizhikevich.github.io/assets/papers/lzr.pdf) appearing at [USENIX Security '21](https://www.usenix.org/conference/usenixsecurity21).

## Building

Install and set up ZMap (see https://github.com/zmap/zmap). If also performing full L7 handshakes, set up ZGrab (https://github.com/zmap/zgrab2).

Set up `$GOPATH` (see https://golang.org/doc/code.html).
```
$ go get github.com/stanford-esrg/lzr
$ cd $GOPATH/src/github.com/stanford-esrg/lzr
```

LZR intercepts connections which ZMap opens; in order to ensure that the kernel does not interfere with LZR, LZR requires a source-ip to be specified for which the kernel drops all RSTs for traffic targeted to the source-ip. The chosen source-ip&mdash;which both ZMap and LZR will use&mdash;should be passed in as a parameter to make, so the appropriate iptables rule can be set.
```
$ make all source-ip=256.256.256.256/32
```

## Usage

To fingerprint unexpected services on an random port (9002):

```
sudo zmap --target-port=9002 --output-filter="success = 1 && repeat = 0" \
-f "saddr,daddr,sport,dport,seqnum,acknum,window" -O json --source-ip=$source-ip | \
sudo ./lzr --handshakes http,tls
```

To complete full L7 handshakes of unexpected services on an random port (9002):

```
sudo zmap --target-port=9002 --output-filter="success = 1 && repeat = 0" \
-f "saddr,daddr,sport,dport,seqnum,acknum,window" -O json --source-ip=$source-ip | \
sudo ./lzr --handshakes wait,http,tls -feedZGrab | \
zgrab multiple -c etc/all.ini 
```

To scan a custom list of IP:Port (i.e., using LZR rather than ZMap to open connections):

```
cat services_list | sudo ./lzr --handshakes http -sendSYNs -sourceIP $source-ip 
```

The expected input format of an example services list is:
```
1.1.1.1:1234
2.2.2.2:80
```


## Flags
```
$ ./lzr --help

Usage of ./lzr:
  -cpuprofile string
    	write cpu profile to file
  -d	debug printing on
  -f string
    	json results output file name (default "default_20210227212802.json")
  -feedZGrab
    	send to zgrab ip and fingerprint
  -forceAllHandshakes
    	Complete all handshakes even if data is returned early on. This also turns off HyperACKtive filtering.
  -gatewayMac string
    	gateway Mac Address in format xx:xx:xx:xx:xx:xx
  -haf int
    	number of random ephemeral probes to send to filter ACKing firewalls
  -handshakes string
    	handshakes to scan with (default "http")
  -memprofile string
    	write memory profile to this file
  -priorityFingerprint string
    	fingerprint to prioritize when multiple match
  -pushDataOnly
    	Don't attach data to ack but rather to push only
  -rn int
    	number of data packets to re-transmit (default 1)
  -rt int
    	number of seconds until re-transmitting packet (default 1)
  -sendInterface string
    	network interface to send packets on (default "ens8")
  -sendSYNs
    	will read input from stdin containing a newline-delimited list of ip:port
  -sourceIP string
    	source IP to send syn packets with (if using sendSYNs flag)
  -t int
    	number of seconds to wait in timeout queue for last retransmission (default 5)
  -w int
    	number of worker threads for each channel (default 1)
```

#### Caveats for specific features
Acking Firewall Filtering (-haf): If a host responds both on the expected port and on the random ephemeral port, whichever response comes first will dictate whether the host is marked as having an ACKing firewall. 

Scanning a custom list of services (-sendSYNs): A sending rate feature has not yet been implemented and therefore all SYNs will be sent at once. Please be careful when using this option to not overload the network. 

## LZR's Algorithm

![](etc/LZRFlow.png)

## License and Copyright

Copyright 2020 The Board of Trustees of The Leland Stanford Junior University

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
