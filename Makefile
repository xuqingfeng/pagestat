HOST_IP := $(shell ifconfig en0 | awk 'FNR == 4 {print $$2}') # https://stackoverflow.com/questions/2382764/escaping-in-makefile
export HOST_IP

tt:
	echo $(HOST_IP)

deps:
	docker-compose up -d

test:
	go test -v $$(go list ./... | grep -v /vendor/)

fmt:
	go fmt ./...

build: fmt
	go build

run: build
	./pagestat
