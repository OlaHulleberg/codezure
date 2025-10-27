GOCACHE := $(PWD)/.gocache

.PHONY: build clean

build:
	GOCACHE="$(GOCACHE)" go build -o bin/codzure .

clean:
	rm -rf bin .gocache
