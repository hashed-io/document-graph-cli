all: build
build:
		mkdir -p bin
		go build -o bin/dgctl cmd/dgctl/main.go
install:
		cd dgctl && go install