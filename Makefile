# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GODEPS=$(GOCMD) get
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt
BINARY_NAME=bin/fake-llm-endpoint
SOURCE_NAME=cmd/fake-llm-endpoint/main.go

all: build

build:
	CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_NAME) -v $(SOURCE_NAME)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

test: build
	$(GOTEST) ./...

fmt:
	$(GOFMT) ./...

deps:
	$(GODEPS) ./...

run: build
	./$(BINARY_NAME)


build-docker: build
	sudo docker build . -t caldito/fake-llm-endpoint:$(VERSION)
