CFLAGS=-Wall

all: cpktgen cmpktgen

cpktgen: cpktgen.o
	gcc -o $@ $<

cmpktgen: cmpktgen.o
	gcc -o $@ $<

.PHONY: clean

clean:
	rm -vf *.o cpktgen cmpktgen
