.PHONY: build run test unit-tests

unit-tests:
	go test -v ./...

run: 
	go run cmd/jarvis.go

build:
	docker build -f Dockerfile -t jarvis:local .

test:
	docker run --rm -it -p 8081:8081 jarvis:local