APP_NAME = denemeler

.PHONY: all build run test clean

all: run

build:
	go build -o bin/$(APP_NAME) main.go

run:
	go run main.go

test:
	go test ./..

clean:
	rm -rf bin/