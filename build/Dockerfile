FROM golang:1.14.4-stretch as builder

COPY go.mod /go/src/github.com/siangyeh8818/prometheus-query-docker/go.mod
COPY go.sum /go/src/github.com/siangyeh8818/prometheus-query-docker/go.sum

# Run golang at any directory, not neccessary $GOROOT, $GOPATH
ENV GO111MODULE=on
WORKDIR /go/src/github.com/siangyeh8818/prometheus-query-docker

# RUN go mod init github.com/pnetwork/sre.monitor.metrics
RUN go mod download
COPY cmd /go/src/github.com/siangyeh8818/prometheus-query-docker/cmd
COPY internal /go/src/github.com/siangyeh8818/prometheus-query-docker/internal
#COPY pkg /go/src/github.com/pnetwork/sre.monitor.metrics/pkg

# Build the Go app
RUN env GOOS=linux GOARCH=amd64 go build -o prometheusQuery -v -ldflags "-s" github.com/siangyeh8818/prometheus-query-docker/cmd/

##### To reduce the final image size, start a new stage with alpine from scratch #####

FROM alpine:3.9
RUN apk --no-cache add bash ca-certificates libc6-compat busybox-extras tzdata
RUN ln -sf /usr/share/zoneinfo/Asia/Taipei /etc/localtime 
RUN echo "Asia/Shanghai" > /etc/timezone
# Run as root
WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /go/src/github.com/siangyeh8818/prometheus-query-docker/prometheusQuery /usr/local/bin/prometheusQuery

# EXPOSE 8081

ENTRYPOINT [ "prometheusQuery" ] 
