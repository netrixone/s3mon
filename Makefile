# Vars and params
GOCMD=go
BINARY_NAME=s3mon

all: build test

clean:
		$(GOCMD) clean -i $(BINARY_NAME)
		rm -f $(BINARY_NAME)

build: deps
		$(GOCMD) install
		$(GOCMD) build -v -o $(BINARY_NAME) github.com/netrixone/s3mon

build-min: deps
		$(GOCMD) install
		$(GOCMD) build -ldflags "-s -w" -v -o $(BINARY_NAME) github.com/netrixone/s3mon

deps:
		$(GOCMD) get -v -t ./...

test: build
		$(GOCMD) test
