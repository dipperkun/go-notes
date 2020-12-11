# grpc-go

## protoc

```bash
wget https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/protoc-3.14.0-linux-x86_64.zip
sudo unzip protoc-3.14.0-linux-x86_64.zip -d /usr
```

## protoc-gen-go & protoc-gen-go-grpc

```bash
go get google.golang.org/protobuf/cmd/protoc-gen-go \
         google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

## hello
```bash
mkdir hello
cd hello
go mod init github.com/dipperkun/go-notes/grpc/hello
mkdir idl
touch idl/greet.proto
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  idl/greet.proto

go get -u google.golang.org/grpc
go mod tidy
```