FROM golang:1.15.6 as golang

RUN apt update
RUN apt install -y gcc
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
RUN export PATH="$PATH:$(go env GOPATH)/bin"

ADD . /PreAmar/
WORKDIR /PreAmar/
RUN go mod tidy
RUN go build -o PreAmar .
