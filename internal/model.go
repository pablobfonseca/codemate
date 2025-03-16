package internal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type Response struct {
	Response string `json:"response"`
}

type StreamResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

type StreamCallback func(chunk string, done bool)

var streamCallback StreamCallback

func SetStreamCallback(callback StreamCallback) {
	streamCallback = callback
}

func SendMessage(prompt, context string) string {
	history := LoadHistory()

	data := Request{
		Model:  "deepseek-coder:6.7b",
		Prompt: fmt.Sprintf("Context: %s\nUser: %s", context+"\n"+history, prompt),
		Stream: true,
	}

	body, _ := json.Marshal(data)
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		errorMsg := fmt.Sprintf("Error connecting to Ollama: %v\n", err)
		if streamCallback != nil {
			streamCallback(errorMsg, true)
		}
		return errorMsg
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	fullResponse := ""

	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			errorMsg := fmt.Sprintf("Error reading stream: %v\n", err)
			if streamCallback != nil {
				streamCallback(errorMsg, true)
			}
			return fullResponse
		}

		var streamResp StreamResponse
		if err := json.Unmarshal(line, &streamResp); err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshaling: %v\n", err)
			continue
		}

		if streamCallback != nil {
			streamCallback(streamResp.Response, streamResp.Done)
		} else {
			fmt.Print(streamResp.Response)
		}

		fullResponse += streamResp.Response

		if streamResp.Done {
			break
		}
	}

	if streamCallback == nil {
		SaveMessage(prompt, fullResponse)
		fmt.Println()
	}

	return fullResponse
}
