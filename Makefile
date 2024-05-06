BIN_NAME := tweet
SRC_PATH := cmd/app/main.go

run:
	go run $(SRC_PATH)

build:
	go build -o bin/$(BIN_NAME) $(SRC_PATH)

install: build
	mv bin/$(BIN_NAME) $(GOBIN)/$(BIN_NAME)

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/$(BIN_NAME)-linux-arm $(SRC_PATH)
	GOOS=linux GOARCH=arm64 go build -o bin/$(BIN_NAME)-linux-arm64 $(SRC_PATH)
	GOOS=freebsd GOARCH=386 go build -o bin/$(BIN_NAME)-freebsd-386 $(SRC_PATH)

all: build install compile
