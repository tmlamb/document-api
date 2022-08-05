FROM golang:1.19.0
RUN go get github.com/go-delve/delve/cmd/dlv
WORKDIR $GOPATH/src
COPY . .
EXPOSE 8080 2345