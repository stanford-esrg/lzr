FROM golang:1.23.0

RUN apt-get update && apt-get install -y libpcap-dev

WORKDIR /go/src/github.com/stanford-esrg/lzr/

COPY . .

RUN go mod tidy \
    && go get -v  gopkg.in/mgo.v2/bson \
    && go get -v github.com/stanford-esrg/lzr \
    && make lzr

CMD ["lzr"]