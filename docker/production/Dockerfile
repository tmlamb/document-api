FROM golangci/golangci-lint AS linter
WORKDIR /go/src/github.com/tmlamb/document-api
COPY . .
RUN golangci-lint run

FROM golang:1.19.0-alpine AS tester
ENV CGO_ENABLED=0
WORKDIR /go/src/github.com/tmlamb/document-api
COPY . .
RUN go test ./... -v -cover

FROM golang:1.15.5-alpine as builder
WORKDIR /go/src/github.com/tmlamb/document-api
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
COPY . .
RUN go build -o build/document-api cmd/document-api/main.go

FROM golang:1.19.0-alpine as runner
RUN apk update && apk add git && go install github.com/jackc/tern@latest
ADD migrations /migrations
COPY --from=builder /go/src/github.com/tmlamb/document-api/build/document-api /usr/bin/document-api
EXPOSE 8080
ENTRYPOINT ["/usr/bin/document-api"]