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

type OpenRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type OpenRouterProvider struct {
	env *config.Env
}

func NewOpenRouterProvider(env *config.Env) *OpenRouterProvider {
	return &OpenRouterProvider{env: env}
}

func (p *OpenRouterProvider) Chat(ctx context.Context, req *ChatRequest) (string, error) {
	apiKey := p.env.OpenRouterAPIKey
	if apiKey == "" {
		return "", fmt.Errorf("openrouter API key is not configured")
	}

	model := req.Model
	if model == "" {
		model = "google/gemini-2.5-flash" // fallback default openrouter model
	}

	messages := []OpenRouterMessage{}
	systemPrompt := ""
	if p.env.ChatbotSystemPrompt != "" {
		if content, err := os.ReadFile(p.env.ChatbotSystemPrompt); err == nil {
			systemPrompt = string(content)
		}
	}
	if systemPrompt != "" {
		messages = append(messages, OpenRouterMessage{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, OpenRouterMessage{Role: "user", Content: req.Prompt})

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

	for k, v := range req.Options {
		payload[k] = v
	}

	bodyData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(bodyData))
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	// OpenRouter-specific headers (optional, but recommended by docs)
	httpReq.Header.Set("HTTP-Referer", "http://localhost:8080")
	httpReq.Header.Set("X-Title", "Marichan API")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		raw, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusRequestEntityTooLarge || resp.StatusCode == 419 || resp.StatusCode == 429 {
			return "", &ProviderError{
				StatusCode: resp.StatusCode,
				Message:    fmt.Sprintf("OpenRouter API limits exceeded: %s", string(raw)),
			}
		}
		return "", fmt.Errorf("openrouter API error: status %d body %s", resp.StatusCode, string(raw))
	}

	var orResp OpenRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&orResp); err != nil {
		return "", err
	}

	if orResp.Error != nil {
		return "", fmt.Errorf("openrouter error: %s", orResp.Error.Message)
	}

	if len(orResp.Choices) > 0 {
		return orResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response choices returned by openrouter")
}
