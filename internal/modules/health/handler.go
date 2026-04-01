package health

import (
	"marichan-api/internal/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Check(c *gin.Context) {
	response.JSON(c, http.StatusOK, true, "ok", gin.H{
		"service": "marichan-api",
	})
}
