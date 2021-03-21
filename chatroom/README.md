# chatroom
A simple chatroom written in Go.

## Build
```shell
git clone https://github.com/dipperkun/go-notes/
cd go-notes/chatroom
go build -o bin/chatroom-server ./server/ # server
go build -o bin/chatroom-cli ./client/ # client
```

## Run
```
./bin/chatroom-server -p 12345 &
# terminal 1
./bin/chatroom-cli -h 127.0.0.1 -p 12345
# terminal 2
./bin/chatroom-cli -h 127.0.0.1 -p 12345
```
