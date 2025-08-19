# fake-llm-endpoint
Fake LLM endpoint following the OpenAI spec. Useful for load testing LLM proxies/gateways without wasting tokens or stressing real servers.

## Features
* OpenAI `v1/chat/completions` compatible.
* Responses have latency like in a real model.
* Streaming and non-streaming responses.
* Supports more than 20k concurrent requests. Benchmarked with k6 and the loadtest.js file available.

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
