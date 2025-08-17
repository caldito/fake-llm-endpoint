package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const maxRequestSize = 1 << 20 // 1MB

// ---- Request Types ----
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

// ---- Response Types ----
type ChatCompletionChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
}

type StreamDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type ChatStreamChunk struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Delta        StreamDelta `json:"delta"`
		Index        int         `json:"index"`
		FinishReason string      `json:"finish_reason,omitempty"`
	} `json:"choices"`
}

// ---- Utility ----
func generateFakeID() string {
	return fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano()+int64(rand.Intn(1000)))
}

// ---- Handler ----
func handleChatCompletions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		http.Error(w, "Expected Content-Type application/json", http.StatusBadRequest)
		return
	}

	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var req ChatCompletionRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	fakeResponse := "This is a fake chat completion."
	processingDelay := time.Duration(12+rand.Intn(7)) * time.Second
	log.Printf("Processing delay: %v", processingDelay)
	time.Sleep(processingDelay)

	id := generateFakeID()
	now := time.Now().Unix()

	if req.Stream {
		w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
			return
		}

		// First delta must include role
		firstChunk := ChatStreamChunk{
			ID:      id,
			Object:  "chat.completion.chunk",
			Created: now,
			Model:   req.Model,
			Choices: []struct {
				Delta        StreamDelta `json:"delta"`
				Index        int         `json:"index"`
				FinishReason string      `json:"finish_reason,omitempty"`
			}{
				{Delta: StreamDelta{Role: "assistant"}, Index: 0},
			},
		}
		data, _ := json.Marshal(firstChunk)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()

		// Stream characters
		for i, ch := range fakeResponse {
			chunk := ChatStreamChunk{
				ID:      id,
				Object:  "chat.completion.chunk",
				Created: now,
				Model:   req.Model,
				Choices: []struct {
					Delta        StreamDelta `json:"delta"`
					Index        int         `json:"index"`
					FinishReason string      `json:"finish_reason,omitempty"`
				}{
					{Delta: StreamDelta{Content: string(ch)}, Index: 0},
				},
			}

			if i == len(fakeResponse)-1 {
				chunk.Choices[0].FinishReason = "stop"
			}

			data, err := json.Marshal(chunk)
			if err != nil {
				log.Printf("Error marshaling stream chunk: %v", err)
				break
			}

			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()

			// Random token delay (50–150 ms)
			time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)
		}

		fmt.Fprint(w, "data: [DONE]\n\n")
		flusher.Flush()
		return
	}

	// Full (non-streaming) response
	resp := ChatCompletionResponse{
		ID:      id,
		Object:  "chat.completion",
		Created: now,
		Model:   req.Model,
		Choices: []ChatCompletionChoice{
			{
				Index: 0,
				Message: ChatMessage{
					Role:    "assistant",
					Content: fakeResponse,
				},
				FinishReason: "stop",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/v1/chat/completions", handleChatCompletions)

	fmt.Println("Listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
