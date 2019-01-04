.PHONY: update install build run test docker
.DEFAULT_GOAL := deploy

novendor = glide novendor | xargs -n 1 go $(1)

update:
	glide update
	go install

install:
	glide install
	go install

build:
	GOOS=linux GOARCH=amd64 go build -ldflags '-extldflags "-static"'

docker:
	docker build --tag=docker.io/glok:latest .

run:
	go run ./main.go -d

test:
	$(call novendor,vet)
	$(call novendor,test)

deploy: install test build docker
