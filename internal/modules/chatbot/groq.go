package chatbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"marichan-api/internal/config"
)

type GroqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
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

type GroqProvider struct {
	env *config.Env
}

func NewGroqProvider(env *config.Env) *GroqProvider {
	return &GroqProvider{env: env}
}

func (p *GroqProvider) Chat(ctx context.Context, req *ChatRequest) (string, error) {
	apiKey := p.env.GroqAPIKey
	if apiKey == "" {
		return "", fmt.Errorf("groq API key is not configured")
	}

	model := req.Model
	if model == "" {
		model = "groq/compound" // fallback default groq model
	}

	messages := []GroqMessage{}
	// For Groq system prompt, we read it just like Gemini did originally if we want to share the prompt,
	// but let's read it here
	systemPrompt := ""
	if p.env.ChatbotSystemPrompt != "" {
		if content, err := os.ReadFile(p.env.ChatbotSystemPrompt); err == nil {
			systemPrompt = string(content)
		}
	}
	if systemPrompt != "" {
		messages = append(messages, GroqMessage{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, GroqMessage{Role: "user", Content: req.Prompt})

	payload := map[string]interface{}{
		"model":    model,
		"messages": messages,
		"stream":   false,
	}

	if req.Temperature != nil {
		payload["temperature"] = *req.Temperature
	}
	if req.TopP != nil {
		payload["top_p"] = *req.TopP
	}
	if req.MaxCompletionTokens > 0 {
		payload["max_completion_tokens"] = req.MaxCompletionTokens
	}

	// Spread custom options directly to Groq payload body
	for k, v := range req.Options {
		payload[k] = v
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
		if resp.StatusCode == 419 {
			return "", &ProviderError{
				StatusCode: resp.StatusCode,
				Message:    "Groq API limits exceeded",
			}
		}
		if resp.StatusCode == http.StatusRequestEntityTooLarge {
			return "", &ProviderError{
				StatusCode: resp.StatusCode,
				Message:    "Request entity too large",
			}
		}
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
