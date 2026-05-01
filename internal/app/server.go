package app

import (
	"fmt"
	"marichan-api/internal/config"
	"marichan-api/internal/database"
	"marichan-api/internal/modules/chatbot"
	"marichan-api/internal/modules/health"
	"marichan-api/internal/modules/pixiv"
	"marichan-api/internal/modules/sticker"
	"marichan-api/internal/modules/traceanime"
	"net/http"
	"path/filepath"

	"gorm.io/gorm"
)

type Server struct {
	Env               *config.Env
	DB                *gorm.DB
	HTTPServer        *http.Server
	HealthHandler     *health.Handler
	StickerHandler    *sticker.Handler
	TraceAnimeHandler *traceanime.Handler
	PixivHandler      *pixiv.Handler
	ChatbotHandler    *chatbot.Handler
}

func NewServer(env *config.Env) (*Server, error) {
	db, err := config.NewDatabase(env)
	if err != nil {
		return nil, err
	}

	if err := database.AutoMigrate(db); err != nil {
		return nil, err
	}

	stickerService := sticker.NewService(env)
	traceAnimeService := traceanime.NewService()
	pixivService := pixiv.NewService(env)
	chatbotService, err := chatbot.NewService(env, db)
	if err != nil {
		return nil, err
	}

	s := &Server{
		Env:               env,
		DB:                db,
		HealthHandler:     health.NewHandler(),
		StickerHandler:    sticker.NewHandler(stickerService),
		TraceAnimeHandler: traceanime.NewHandler(traceAnimeService),
		PixivHandler:      pixiv.NewHandler(pixivService),
		ChatbotHandler:    chatbot.NewHandler(chatbotService),
	}

	router := setupRouter(s)

	s.HTTPServer = &http.Server{
		Addr:    ":" + env.AppPort,
		Handler: router,
	}

	absUpload, _ := filepath.Abs(env.UploadDir)
	fmt.Println("upload dir:", absUpload)

	return s, nil
}

func (s *Server) Run() error {
	return s.HTTPServer.ListenAndServe()
}
