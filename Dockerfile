# Build the manager binary
FROM golang:1.15-alpine as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY ./ ./

# Build
RUN go build -o=./initializer ./main.go

FROM centos:7.3.1611
ARG HOME=/home/work
WORKDIR ${HOME}
COPY --from=builder /workspace/initializer ${HOME}
CMD ["./initializer"]
