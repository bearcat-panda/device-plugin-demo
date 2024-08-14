FROM registry.cn-hangzhou.aliyuncs.com/bearcat-panda/golang:1.22 AS builder

WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
ENV GOPROXY=https://goproxy.cn,direct


# Copy the entire project
COPY . .
RUN go mod tidy

# Build the project
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/device-plugin cmd/main.go

FROM registry.cn-hangzhou.aliyuncs.com/bearcat-panda/alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/bin/device-plugin .

ENTRYPOINT ["./device-plugin"]
