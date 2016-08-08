all: fmt get build

fmt:
	go fmt ./...

get:
	go get ./...

build:
	go build -o target/google-calendar-retreiver ./...
