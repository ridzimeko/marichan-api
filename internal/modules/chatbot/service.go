package chatbot

import (
	"context"
	"fmt"
	"marichan-api/internal/config"
)

type ProviderError struct {
	StatusCode int
	Message    string
}

func (e *ProviderError) Error() string {
	return e.Message
}

type Provider interface {
	Chat(ctx context.Context, req *ChatRequest) (string, error)
}

type Service struct {
	providers map[string]Provider
}

func NewService(env *config.Env) (*Service, error) {
	geminiProv, err := NewGeminiProvider(env)
	if err != nil {
		return nil, err
	}
	groqProv := NewGroqProvider(env)
	openrouterProv := NewOpenRouterProvider(env)

	return &Service{
		providers: map[string]Provider{
			"gemini":     geminiProv,
			"groq":       groqProv,
			"openrouter": openrouterProv,
		},
	}, nil
}

func (s *Service) Chat(ctx context.Context, req *ChatRequest) (string, error) {
	providerName := req.Provider
	if providerName == "" {
		providerName = "gemini"
	}

	prov, ok := s.providers[providerName]
	if !ok {
		return "", fmt.Errorf("unknown chatbot provider: %s", providerName)
	}

	return prov.Chat(ctx, req)
}
