# Compile source code.
FROM golang:1.14-stretch as builder
ENV SRC_DIR="/go/src/github.com/satriajidam/go-gin-skeleton"
WORKDIR $SRC_DIR
COPY go.mod .
COPY go.sum .
# Copy all source and build it.
# This layer will be rebuilt whenever a file has changed in the source directory.
COPY ./ .
RUN GOOS=linux GOPROXY=https://proxy.golang.org go build -v -a -mod=readonly -o /bin/server .

# Build final image.
FROM debian:stretch-slim
RUN apt-get update -y \
  && apt-get install -y --no-install-recommends \
    ca-certificates \
    openssl \
    curl \
  && update-ca-certificates \
  && apt-get clean && rm -rf /tmp/* /var/tmp/* /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /bin/server server
ENTRYPOINT ["./server"]
