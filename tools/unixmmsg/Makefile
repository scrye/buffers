CFLAGS=-Wall

all: client server

client: client.o
	gcc -o $@ $<

server: server.o
	gcc -o $@ $<

.PHONY: clean

clean:
	rm -vf *.o client server
