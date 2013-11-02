GOPATH := ${PWD}
export GOPATH

all: server

server: server.go libassets.a
	go build server.go

libassets.a: assets.o
	ar -r libassets.a assets.o

assets.o: assets.zip
	objcopy -B i386 -I binary -O elf64-x86-64 assets.zip assets.o

assets.zip: 
	zip -r assets.zip assets

clean:
	rm assets.zip assets.o libassets.a server
