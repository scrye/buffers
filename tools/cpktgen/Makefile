CFLAGS=-Wall

all: cpktgen

cpktgen: cpktgen.o
	gcc -o $@ $<

.PHONY: clean

clean:
	rm -vf *.o cpktgen
