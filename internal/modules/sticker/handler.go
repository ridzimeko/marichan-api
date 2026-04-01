package sticker

import (
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

func (h *Handler) Convert(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "image is required", gin.H{
			"image": []string{"file is required"},
		})
		return
	}

	result, err := h.service.ConvertToSticker(c.Request.Context(), file)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.JSON(c, http.StatusOK, true, "sticker created successfully", result)
}
