package chatbot

import (
	"context"
	"fmt"
	"strings"

	"marichan-api/internal/config"
	"marichan-api/internal/database"

	"github.com/pgvector/pgvector-go"
	"google.golang.org/genai"
	"gorm.io/gorm"
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
	env       *config.Env
	db        *gorm.DB
	providers map[string]Provider
}

func NewService(env *config.Env, db *gorm.DB) (*Service, error) {
	geminiProv, err := NewGeminiProvider(env)
	if err != nil {
		return nil, err
	}
	groqProv := NewGroqProvider(env)
	openrouterProv := NewOpenRouterProvider(env)

	return &Service{
		env: env,
		db:  db,
		providers: map[string]Provider{
			"gemini":     geminiProv,
			"groq":       groqProv,
			"openrouter": openrouterProv,
		},
	}, nil
}

func (s *Service) getEmbedding(ctx context.Context, text string) ([]float32, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: s.env.GeminiAPIKey})
	if err != nil {
		return nil, err
	}

	contents := []*genai.Content{{Parts: []*genai.Part{{Text: text}}}}
	dim := int32(768)
	resp, err := client.Models.EmbedContent(ctx, "gemini-embedding-2", contents, &genai.EmbedContentConfig{
		TaskType:             "SEMANTIC_SIMILARITY",
		OutputDimensionality: &dim,
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return resp.Embeddings[0].Values, nil
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

	var userEmbedding []float32
	var originalPrompt = req.Prompt

	// If sessionID is provided, use Vector DB RAG
	if req.SessionID != "" && s.db != nil {
		emb, err := s.getEmbedding(ctx, originalPrompt)
		if err == nil && len(emb) > 0 {
			userEmbedding = emb

			// Find nearest 5 past conversations
			var pastChats []database.ChatMemory
			s.db.Order("embedding <-> '"+pgvector.NewVector(emb).String()+"'").
				Where("session_id = ?", req.SessionID).
				Limit(5).
				Find(&pastChats)

			if len(pastChats) > 0 {
				var contextBuilder strings.Builder
				contextBuilder.WriteString("The following are relevant past messages from this conversation history to provide context:\n")
				for _, chat := range pastChats {
					contextBuilder.WriteString(fmt.Sprintf("[%s]: %s\n", chat.Role, chat.Content))
				}
				contextBuilder.WriteString("\nBased on the context above, please answer the new prompt.\n\nNew prompt: ")
				contextBuilder.WriteString(originalPrompt)

				req.Prompt = contextBuilder.String()
			}
		} else {
			fmt.Println("Warning: Failed to get embedding for RAG:", err)
		}
	}

	reply, err := prov.Chat(ctx, req)
	if err != nil {
		return "", err
	}

	// Restore original prompt so we save only the prompt, not the injected context
	req.Prompt = originalPrompt

	// Save to DB asynchronously
	if req.SessionID != "" && s.db != nil && len(userEmbedding) > 0 {
		go func(sessionID, prompt, reply string, uEmb []float32) {
			bgCtx := context.Background()

			// Save User Prompt
			s.db.Create(&database.ChatMemory{
				SessionID: sessionID,
				Role:      "user",
				Content:   prompt,
				Embedding: pgvector.NewVector(uEmb),
			})

			// Save Assistant Reply
			replyEmb, err := s.getEmbedding(bgCtx, reply)
			if err == nil {
				s.db.Create(&database.ChatMemory{
					SessionID: sessionID,
					Role:      "assistant",
					Content:   reply,
					Embedding: pgvector.NewVector(replyEmb),
				})
			}
		}(req.SessionID, originalPrompt, reply, userEmbedding)
	}

	return reply, nil
}
