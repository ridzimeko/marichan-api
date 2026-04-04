package app

import (
	"marichan-api/internal/middleware"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setupRouter(server *Server) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	r.MaxMultipartMemory = server.Env.MaxUploadSizeMB << 20

	r.Static("/files", server.Env.UploadDir)

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "marichan-api running",
		})
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", server.HealthHandler.Check)
		v1.GET("/pixiv/download-artwork", server.PixivHandler.DownloadArtwork)

		protected := v1.Group("/")
		protected.Use(middleware.APIKeyAuth(server.Env))
		{
			protected.POST("/sticker/convert", server.StickerHandler.Convert)
			protected.POST("/trace-anime/search", server.TraceAnimeHandler.Search)

			// pixiv routes
			protected.GET("/pixiv/artworks/:id", server.PixivHandler.GetArtworkDetail)
			protected.GET("/pixiv/search/artworks", server.PixivHandler.SearchArtworks)
			protected.GET("/pixiv/artists/:id", server.PixivHandler.GetArtistDetail)
			protected.GET("/pixiv/download", server.PixivHandler.DownloadImage)
			protected.POST("/chatbot/chat", server.ChatbotHandler.Chat)
		}
	}

	return r
}
