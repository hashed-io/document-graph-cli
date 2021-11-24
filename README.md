### Command line application for accessing Document Graph

# Quick Start 
## Prerequisites
- Go 1.16+
- cmake

## Build
```
git clone git@github.com:hashed-io/document-graph-cli.git
cd document-graph-cli
git checkout develop
make
./bin/dgctl
```
## Menu
```
    Command line application written in Go for 
        interacting with an on-chain Document Graph

Usage:
  dgctl [command]

Available Commands:
  create      create documents or edges in the Document Graph
  deploy      deploy contracts
  gen         generate graph data
  get         get objects from the Document Graph
  help        Help about any command

Flags:
  -c, --config string       configuration file (default ".dgctl.yaml")
      --expiration int      Set time before transaction expires, in seconds. Defaults to 30 seconds. (default 30)
  -h, --help                help for dgctl
      --vault-file string   Wallet file that contains encrypted key material (default "./eosc-vault.json")

Use "dgctl [command] --help" for more information about a command.
```


## Using insecure vault (DEV/TEST only)
You can set Vault to false and have the plain text keys read. You will then not be prompted for to use a secure vault.

```yaml
EosioEndpoint: https://test.telos.kitchen
AssetsAsFloat: true
Contract: hasheddocgrh
UserAccount: hasheduser11
GeneratedDir: _generated
Vault: false
Keys: ["5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3", "5Kdp6CQRq6MwZVjCFcfjSNeGceD3RZ1rtcQiyaz7sEv7SdR4E6r"]
```

# Useful Tips

Use ```DEBUG=true``` before command to show verbose debug messages.
```bash
DEBUG=true ./bin/dgctl --config .testnet.yaml get documents
```

To run development/local updates
```bash
go run cmd/dgctl/main.go --config .testnet.yaml get documents
```

Create some random authors, posts, and likes
```
./bin/dgctl --config .testnet.yaml gen --users 10 --posts 50
```

View a specific document and begin navigating
```
./bin/dgctl --config .testnet.yaml get document 35
```

# Generated Social Media data structure
![generated-social-media-data-structure](//www.plantuml.com/plantuml/png/SoWkIImgAStDuKh9J2zABCXGS5Uevb80WgB4lEoKp29Rdo1hC3-qEBL8GTTE8I2_k4GXEYSnAJN7LYcnj2GZloWrHIaMZwASp6mC5M0QhY8jFoSdlxmOfAQMoo4rBmNe8000)