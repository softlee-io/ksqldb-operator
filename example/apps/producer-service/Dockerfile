FROM golang:latest as build

RUN apt-get install gcc

ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

RUN mkdir -p /workdir
ADD . /workdir
WORKDIR /workdir

RUN go mod download
RUN go build -o ./server src/main.go

CMD ["./server"]
EXPOSE 8080