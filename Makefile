IMG = registry.cn-hangzhou.aliyuncs.com/bearcat-panda/device-plugin:latest

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -o bin/device-plugin cmd/main.go

.PHONY: build-image
build-image:
	docker build -t ${IMG} .