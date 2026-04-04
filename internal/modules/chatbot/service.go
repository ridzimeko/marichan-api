package chatbot

import (
	"context"
	"marichan-api/internal/config"
	"os"

	"google.golang.org/genai"
)

type ProviderError struct {
	StatusCode int
	Message    string
}

func (e *ProviderError) Error() string {
	return e.Message
}

type Service struct {
	env    *config.Env
	client *genai.Client
}

func NewService(env *config.Env) (*Service, error) {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{APIKey: env.GeminiAPIKey})
	if err != nil {
		return nil, err
	}
	return &Service{
		env:    env,
		client: client,
	}, nil
}

func (s *Service) Chat(ctx context.Context, req *ChatRequest) (string, error) {
	systemPrompt := ""
	if s.env.GeminiSystemPromptFile != "" {
		if content, err := os.ReadFile(s.env.GeminiSystemPromptFile); err == nil {
			systemPrompt = string(content)
		}
	}

	if req.Provider == "groq" {
		return s.chatGroq(ctx, s.env.GroqAPIKey, systemPrompt, req)
	}

	// Default fallback is Gemini
	config := &genai.GenerateContentConfig{}
	
	if systemPrompt != "" {
		config.SystemInstruction = &genai.Content{
			Parts: []*genai.Part{{Text: systemPrompt}},
		}
	}

	if req.Temperature != nil {
		config.Temperature = req.Temperature
	}
	if req.TopP != nil {
		config.TopP = req.TopP
	}
	if req.MaxCompletionTokens > 0 {
		config.MaxOutputTokens = int32(req.MaxCompletionTokens)
	}

	resp, err := s.client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(req.Prompt), config)
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) > 0 {
		for _, part := range resp.Candidates[0].Content.Parts {
			if part.Text != "" {
				return part.Text, nil
			}
		}
	}

	return "", nil
}
