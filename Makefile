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


build-docker: build # This is for local testing only, images are built and pushed with gh actions
	sudo docker build . -t ghcr.io/caldito/fake-llm-endpoint:$(VERSION)

load-test:
	sudo docker run --rm -i --network="host" -e BASE_URL=http://localhost:8080 -v $(pwd)/loadtest.js:/loadtest.js grafana/k6 run /loadtest.js
