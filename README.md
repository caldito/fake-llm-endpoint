# fake-llm-endpoint
Fake LLM endpoint following the OpenAI spec. Usefull for load testing LLM proxies/gateways and cheap but real integration tests.

## Features
* OpenAI `v1/chat/completions` compatible.
* Responses have latency like in a real model.
* Streaming and non-streaming responses.
* Supports a high number of concurrent requests. (To be benchmarked)

## Run with docker
Start the server:

```
docker run -p 8080:8080 -d ghcr.io/caldito/fake-llm-endpoint:0.1.0
```

## Development
Uses standard Go libraries. Running on go 1.25.0.

to build the binary run
```
make
```

to build and run the binary introduce
```
make run
```
