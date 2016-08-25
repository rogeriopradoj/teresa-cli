FROM golang:1.6

RUN mkdir -p /go/src/github.com/luizalabs/tcli
WORKDIR /go/src/github.com/luizalabs/tcli
COPY . /go/src/github.com/luizalabs/tcli

RUN go get github.com/tools/godep
RUN godep go install .

CMD ["tcli"]
