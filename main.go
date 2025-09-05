package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"webhook/internal/handler"
	"webhook/internal/infrastructure/config"
	"webhook/internal/usecase"
)

func main() {
	// 設定 Zerolog，使用適合開發的 ConsoleWriter
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// 初始化配置
	cfg := config.NewConfig()

	// 建立 Gin 引擎
	r := gin.Default()

	// UseCase
	jiraConfig := cfg.GetJiraConfig()
	execClaudeCode := usecase.NewExecClaudeCode(log.Logger, jiraConfig)

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
