package handler

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"webhook/internal/handler/dto"
	"webhook/internal/usecase"
)

type Handlers struct {
	execClaudeCode *usecase.ExecClaudeCode
}

// NewHandler 建立新的 Handlers
func NewHandler(ExecClaudeCode *usecase.ExecClaudeCode) *Handlers {
	return &Handlers{
		execClaudeCode: ExecClaudeCode,
	}
}

// WebhookHandler 處理所有來自 Notion 的請求
func (h *Handlers) WebhookHandler(c *gin.Context) {
	// 讀取並記錄 request body
	_, err := logAndRestoreRequestBody(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read request"})
		return
	}

	var payload dto.NotionWebhookDto

	// 將傳入的 JSON 綁定到結構體
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Error().Err(err).Msg("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// print payload log
	log.Info().Any("payload", payload).Msg("Received webhook payload")

	if payload.VerificationToken != "" {
		log.Info().Str("verification_token", payload.VerificationToken).Msg("Received verification token")
		c.Status(http.StatusOK)
		return
	}

	if !payload.IsAssignUpdate() {
		c.Status(http.StatusOK)
		return
	}

	go func() {
		ticketId := payload.Entity.ID
		log.Info().Str("pmTool", "notion").Str("ticket_id", ticketId).Msg("收到任務")

		if execErr := h.execClaudeCode.ExecTask(ticketId, "notion"); execErr != nil {
			log.Error().
				Str("ticketId", ticketId).
				Err(execErr).
				Str("pmTool", "notion").
				Msg("Failed to execClaudeCode")
		}

	}()

	c.Status(http.StatusOK)
}

// HealthHandler 處理健康檢查請求
func (h *Handlers) HealthHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

// MyTestHandler 測試用Handler //TODO: 刪除
func (h *Handlers) MyTestHandler(c *gin.Context) {

	// TODO: thread pool
	ticketId := "2575ab0c0f4c80db89c5ebace1df2173"

	log.Info().Str("ticket_id", ticketId).Msg("收到任務")

	// cmd := exec.Command("claude",
	// 	"--dangerously-skip-permissions",
	// 	"-p", "'產生一個包含hi aaa 的txt file'")

	// output, err := cmd.Output()
	// if err != nil {
	// 	fmt.Printf("錯誤: %v\n", err.Error())
	// 	return
	// }
	// fmt.Printf("輸出:\n%s", output)

	go func() {
		if execErr := h.execClaudeCode.ExecTask(ticketId, "notion"); execErr != nil {
			log.Error().
				Str("ticketId", ticketId).
				Err(execErr).
				Str("pmTool", "notion").
				Msg("Failed to execClaudeCode")
		}

	}()

	c.String(http.StatusOK, "OK")

}

// RegisterRoutes 註冊所有路由
func (h *Handlers) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.HealthHandler)
	r.POST("/notion/webhook", h.WebhookHandler) // webhook路徑
	r.POST("/mytest", h.MyTestHandler)
}

// logAndRestoreRequestBody 讀取並印出 request body，然後重新設定給 context
func logAndRestoreRequestBody(c *gin.Context) ([]byte, error) {
	// 讀取原始 request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read request body")
		return nil, err
	}

	// 印出原始 JSON
	log.Info().RawJSON("raw_payload", body).Msg("Received raw webhook payload")

	// 將 body 放回去，讓後續的 ShouldBindJSON 可以使用
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	return body, nil
}
