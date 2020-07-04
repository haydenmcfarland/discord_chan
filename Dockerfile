FROM golang:alpine

COPY . /go/src/github.com/haydenmcfarland/discord_chan
WORKDIR /go/src/github.com/haydenmcfarland/discord_chan

RUN go get ./
RUN go build

CMD discord_chan
