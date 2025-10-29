GOCACHE := $(PWD)/.gocache

.PHONY: build clean

build:
    GOCACHE="$(GOCACHE)" go build -o bin/codezure .

clean:
	rm -rf bin .gocache
