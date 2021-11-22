### Command line application for accessing Document Graph

# Quick Start 
## Prerequisites
- Go 1.16+
- cmake

## Build
```
make
./bin/dgctl
```


# Useful Commands

Print the documents on the testnet document graph
```
DEBUG=true go run cmd/dgctl/main.go --config .testnet.yaml get documents
```

Create some random authors, posts, and likes
```
DEBUG=true go run cmd/dgctl/main.go --config .testnet.yaml gen --users 10 --posts 50
```

View a specific document and begin navigating
```
DEBUG=true go run cmd/dgctl/main.go --config .testnet.yaml get document 35
```