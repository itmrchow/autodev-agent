package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"webhook/internal/handler"
	"webhook/internal/usecase"
)

func main() {
	// 設定 Zerolog，使用適合開發的 ConsoleWriter
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// 建立 Gin 引擎
	r := gin.Default()

	// UseCase
	execClaudeCode := usecase.NewExecClaudeCode(log.Logger)

	// Handlers
	handlers := handler.NewHandler(execClaudeCode)

	// 註冊路由
	handlers.RegisterRoutes(r)

	// 啟動伺服器
	log.Info().Msg("Starting server on port :8090")
	if err := r.Run(":8090"); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
