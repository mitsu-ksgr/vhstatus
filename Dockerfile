#------------------------------------------------
# Develop
FROM golang:1.16.2-alpine AS Development
RUN apk --no-cache add upx build-base

WORKDIR /go/src/github.com/mitsu-ksgr/vhstatus

CMD ["go", "run", "cmd/main.go", "-port", "8000", \
  "-log-dir-path", "/go/src/github.com/mitsu-ksgr/vhstatus/test/logs", \
  "-template-dir-path", "/go/src/github.com/mitsu-ksgr/vhstatus/web"]


#------------------------------------------------
# Builder
FROM golang:1.16.2-alpine AS Builder
RUN apk --no-cache add upx

WORKDIR /go/src/github.com/mitsu-ksgr/vhstatus

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -ldflags="-w" -o vhstatus-server cmd/main.go
RUN upx vhstatus-server


#------------------------------------------------
# binary to local
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /go/src/github.com/mitsu-ksgr/vhstatus

COPY --from=builder /go/src/github.com/mitsu-ksgr/vhstatus/vhstatus-server /root/vhstatus-server

CMD ["cp", "/root/vhstatus-server", "./"]

