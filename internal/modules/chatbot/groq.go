package chatbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GroqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqRequest struct {
	Model           string        `json:"model"`
	Messages        []GroqMessage `json:"messages"`
	Temperature     float32       `json:"temperature,omitempty"`
	MaxTokens       int           `json:"max_completion_tokens,omitempty"`
	TopP            float32       `json:"top_p,omitempty"`
	Stream          bool          `json:"stream"`
	ReasoningEffort string        `json:"reasoning_effort,omitempty"`
}

type GroqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (s *Service) chatGroq(ctx context.Context, apiKey string, systemPrompt string, req *ChatRequest) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("groq API key is not configured")
	}

	model := req.Model
	if model == "" {
		model = "openai/gpt-oss-120b" // fallback default groq model
	}

	messages := []GroqMessage{}
	if systemPrompt != "" {
		messages = append(messages, GroqMessage{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, GroqMessage{Role: "user", Content: req.Prompt})

	temp := float32(1.0)
	if req.Temperature != nil {
		temp = *req.Temperature
	}

	topP := float32(1.0)
	if req.TopP != nil {
		topP = *req.TopP
	}

	payload := GroqRequest{
		Model:       model,
		Messages:    messages,
		Temperature: temp,
		TopP:        topP,
		MaxTokens:   req.MaxCompletionTokens,
		Stream:      false,
	}

	bodyData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewReader(bodyData))
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		raw, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("groq API error: status %d body %s", resp.StatusCode, string(raw))
	}

	var groqResp GroqResponse
	if err := json.NewDecoder(resp.Body).Decode(&groqResp); err != nil {
		return "", err
	}

	if groqResp.Error != nil {
		return "", fmt.Errorf("groq error: %s", groqResp.Error.Message)
	}

	if len(groqResp.Choices) > 0 {
		return groqResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response choices returned by groq")
}
