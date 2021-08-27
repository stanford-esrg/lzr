FROM golang:1.14

RUN apt-get update && apt-get install -y libpcap-dev

RUN go get -v gopkg.in/mgo.v2/bson
RUN go get -v github.com/stanford-esrg/lzr

COPY . /go/src/github.com/stanford-esrg/lzr/

RUN (cd /go/src/github.com/stanford-esrg/lzr && make lzr)

WORKDIR /go/src/github.com/stanford-esrg/lzr

CMD ["lzr"]