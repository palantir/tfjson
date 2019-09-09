# Build 'tfjson' executable binary
FROM golang:alpine AS builder

RUN apk add --no-cache git
RUN adduser -D -g '' tfjson

WORKDIR $GOPATH/src/palantir/tfjson/
COPY . .

RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/tfjson

# Build small image only containing 'tfjson' binary
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/tfjson /go/bin/tfjson

USER tfjson

ENTRYPOINT ["/go/bin/tfjson"]
