FROM golangci/golangci-lint AS linter
WORKDIR /go/src/github.com/tmlamb/document-api
COPY . .
RUN golangci-lint run

FROM golang:1.19.0-alpine AS tester
ENV CGO_ENABLED=0
WORKDIR /go/src/github.com/tmlamb/document-api
COPY . .
RUN go test ./... -v -cover

FROM golang
WORKDIR /go/src/github.com/tmlamb/document-api
ADD . /go/src/github.com/tmlamb/document-api
RUN go get -v ./...
RUN go install github.com/jackc/tern@latest