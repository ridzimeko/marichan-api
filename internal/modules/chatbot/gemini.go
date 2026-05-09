package chatbot

import (
	"context"
	"encoding/base64"
	"marichan-api/internal/config"
	"os"

	"google.golang.org/genai"
)

type GeminiProvider struct {
	env    *config.Env
	client *genai.Client
}

func NewGeminiProvider(env *config.Env) (*GeminiProvider, error) {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{APIKey: env.GeminiAPIKey})
	if err != nil {
		return nil, err
	}
	return &GeminiProvider{
		env:    env,
		client: client,
	}, nil
}

func (g *GeminiProvider) Chat(ctx context.Context, req *ChatRequest) (string, error) {
	systemPrompt := ""
	if g.env.ChatbotSystemPrompt != "" {
		if content, err := os.ReadFile(g.env.ChatbotSystemPrompt); err == nil {
			systemPrompt = string(content)
		}
	}

	var tools = []*genai.Tool{
		{
			GoogleSearch: &genai.GoogleSearch{},
		},
	}

	var config *genai.GenerateContentConfig = &genai.GenerateContentConfig{
		Tools: tools,
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingLevel: genai.ThinkingLevelMinimal,
		},
	}

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

	var parts []*genai.Part
	for _, f := range req.Files {
		if f.Data != "" {
			decoded, err := base64.StdEncoding.DecodeString(f.Data)
			if err == nil {
				parts = append(parts, &genai.Part{
					InlineData: &genai.Blob{
						MIMEType: f.MimeType,
						Data:     decoded,
					},
				})
			}
		} else if f.URI != "" {
			parts = append(parts, &genai.Part{
				FileData: &genai.FileData{
					FileURI:  f.URI,
					MIMEType: f.MimeType,
				},
			})
		}
	}

	parts = append(parts, &genai.Part{
		Text: req.Prompt,
	})

	contents := []*genai.Content{
		{
			Role:  "user",
			Parts: parts,
		},
	}

	resp, err := g.client.Models.GenerateContent(ctx, "gemma-4-26b-a4b-it", contents, config)
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
