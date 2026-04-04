package chatbot

import (
	"errors"
	"marichan-api/internal/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type ChatRequest struct {
	Prompt              string   `json:"prompt"`
	Provider            string   `json:"provider"` // e.g. "gemini" or "groq"
	Model               string   `json:"model"`
	Temperature         *float32 `json:"temperature"`
	MaxCompletionTokens int      `json:"max_completion_tokens"`
	TopP                *float32 `json:"top_p"`
}

func (h *Handler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request format", nil)
		return
	}
	if req.Prompt == "" {
		response.Error(c, http.StatusBadRequest, "prompt is required", nil)
		return
	}

	reply, err := h.service.Chat(c.Request.Context(), &req)
	if err != nil {
		var provErr *ProviderError
		if errors.As(err, &provErr) {
			response.Error(c, provErr.StatusCode, provErr.Message, nil)
			return
		}
		response.Error(c, http.StatusBadGateway, err.Error(), nil)
		return
	}

	response.JSON(c, http.StatusOK, true, "Chat successful", gin.H{"reply": reply})
}
