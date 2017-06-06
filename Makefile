deps:
	docker-compose up -d

test:
	go test -v $$(go list ./... | grep -v /vendor/)

fmt:
	go fmt ./...

build: fmt
	go build ./cmd/pagestat -o ./out

run: build
	./out/pagestat -mode broker
