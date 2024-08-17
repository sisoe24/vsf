all: test build

BIN=vsf
BUILD=./build/$(BIN)

build:
	go build -o $(BUILD) -v -ldflags "-w -s"

test:
	go test -v ./...

clean:
	go clean
	rm -f $(BUILD)

run: build
	$(BUILD)

install:
	go install

lint:
	staticcheck ./...

.PHONY: all build test clean run deps install lint
