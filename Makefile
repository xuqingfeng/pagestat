deps:
	docker-compose up -d

test: deps
	go test -v $$(go list ./... | grep -v /vendor/)

fmt:
	go fmt ./...

build: fmt
	go build -o ./out/pagestat ./cmd/pagestat

run: build
	./out/pagestat -mode broker
