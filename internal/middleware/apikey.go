package middleware

import (
	"marichan-api/internal/config"
	"marichan-api/internal/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth(env *config.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-MARICHAN-API-KEY")
		if apiKey == "" || apiKey != env.APIKey {
			response.Error(ctx, http.StatusUnauthorized, "unauthorized", gin.H{
				"api_key": []string{"Invalid or missing API key"},
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
