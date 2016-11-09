FROM golang:1.6

RUN mkdir -p /go/src/github.com/luizalabs/teresa-cli
WORKDIR /go/src/github.com/luizalabs/teresa-cli
COPY . /go/src/github.com/luizalabs/teresa-cli

RUN go get github.com/tools/godep
RUN godep go install .

CMD ["tcli"]
