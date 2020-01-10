.PHONY: build run test docker
.DEFAULT_GOAL := deploy


install:
	glide install
	go install

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-extldflags "-static"'

docker:
	docker build --tag=docker.io/ondrejit/glok:latest .

run:
	go run ./main.go -d

test:
	go mod verify
	go test ./...

deploy: test build docker
