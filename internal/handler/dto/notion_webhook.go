package dto

import (
	"strings"
	"time"
)

var _ PmToolReqInterface = &NotionWebhookDto{}

// NotionWebhookDto 是 Notion webhook 事件的頂級結構體
type NotionWebhookDto struct {
	APIVersion     string         `json:"api_version"`
	AttemptNumber  int            `json:"attempt_number"`
	Authors        []NotionAuthor `json:"authors"`
	Data           *NotionData    `json:"data,omitempty"`
	Entity         NotionEntity   `json:"entity"`
	ID             string         `json:"id"`
	IntegrationID  string         `json:"integration_id"`
	SubscriptionID string         `json:"subscription_id"`
	Timestamp      time.Time      `json:"timestamp"`
	Type           string         `json:"type"`
	WorkspaceID    string         `json:"workspace_id"`
	WorkspaceName  string         `json:"workspace_name"`

	VerificationToken string `json:"verification_token"`
}

// NotionAuthor 代表作者資訊
type NotionAuthor struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// NotionData 包含 webhook 的數據部分
type NotionData struct {
	Parent            NotionParent `json:"parent"`
	UpdatedProperties []string     `json:"updated_properties"`
}

// NotionParent 代表父級對象
type NotionParent struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// NotionEntity 代表實體對象
type NotionEntity struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func (dto *NotionWebhookDto) IsAssignUpdate() bool {
	// 檢查是否為頁面屬性更新事件
	if dto.Type != "page.properties_updated" {
		return false
	}

	// 檢查是否有數據
	if dto.Data == nil {
		return false
	}

	// 檢查是否有 ID
	if dto.Entity.ID == "" {
		return false
	}

	// 檢查更新的屬性中是否包含 assign_property
	for _, property := range dto.Data.UpdatedProperties {
		if strings.Contains(property, "assign_property") {
			return true
		}
	}

	return false
}
