package chatbot

import (
	"context"
	"marichan-api/internal/config"
	"os"

	"google.golang.org/genai"
)

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
	var config *genai.GenerateContentConfig
	if systemPrompt != "" {
		config = &genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{{Text: systemPrompt}},
			},
		}
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
