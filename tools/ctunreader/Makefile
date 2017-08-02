CFLAGS=-Wall

all: ctunreader

ctunreader: ctunreader.o
	gcc -o $@ $<

.PHONY: clean

clean:
	rm -vf *.o ctunreader
