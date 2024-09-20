# Debian 11

FROM golang:1.22.6-bullseye

RUN apt-get update && \
    apt-get install -y build-essential

# Go tools

RUN go install golang.org/x/tools/cmd/goyacc@latest

RUN go install -v github.com/go-delve/delve/cmd/dlv@latest

# Switch to a working directory

WORKDIR /workspace