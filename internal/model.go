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
		return fmt.Sprintf("Error: %v\n", err)
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
			fmt.Fprintf(os.Stderr, "Error reading stream: %v\n", err)
			break
		}

		var streamResp StreamResponse
		if err := json.Unmarshal(line, &streamResp); err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshaling: %v\n", err)
			continue
		}

		fmt.Print(streamResp.Response)
		fullResponse += streamResp.Response

		if streamResp.Done {
			break
		}
	}

	SaveMessage(prompt, fullResponse)

	fmt.Println()
	return fullResponse
}
