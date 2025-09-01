# autodev-agent

這個服務是一個 SPEC 驅動的 webhook 處理器，負責接收來自專案管理工具的事件，並觸發 AI Agent 進行任務  

希望可以改善與實踐  
- Spec 驅動的開發
- 快速產出 prototype , 供 RD review , 加速產出
- Claude Code (or other AI tool) 並行開發

## Features
- [x] 接收 PM 工具事件
  - [x] Notion
  - [ ] Jira
  - Other...

- [x] 觸發Agent開發 , 更新 ticket

## tech stack
- go + gin : webhook
- notion : pm tool
- claude : ai agent

## 部署與設定

### Prerequisites
- go 1.24+
- Docker (for containerized deployment)
- webhook knowledge

### Environment Variables

#### 核心服務設定
<!-- - `PM_TOOL_NOTION_TOKEN` Notion API token
- `PM_TOOL_JIRA_TOKEN` Jira API token (optional) -->
- `CODE_REPO_GITHUB_TOKEN` GitHub API token
- `PORT` 服務端口，預設 `8090`

#### GitHub 配置（用於 Claude Code 開發功能）
- `GITHUB_USER_NAME` GitHub 使用者名稱
- `GITHUB_USER_EMAIL` GitHub 使用者 email
- `GITHUB_PERSONAL_TOKEN` GitHub Personal Access Token

<!-- #### Claude Code 配置
- `CLAUDE_API_KEY` Claude API key（用於 AI 開發功能） -->

### 部署

#### Use Docker image

```bash
# 拉取最新映像
docker pull itmrchow/autodev-agent:v0.1.0

# 運行容器（包含完整配置）
docker run -d \
  -p 8090:8090 \
  -e CODE_REPO_GITHUB_TOKEN=your_github_token \
  -e GITHUB_USER_NAME="Your Name" \
  -e GITHUB_USER_EMAIL="your.email@example.com" \
  -e GITHUB_PERSONAL_TOKEN=your_github_personal_token \
  -v $(pwd)/projects:/app/projects \
  --name autodev-agent \
  itmrchow/autodev-agent:v0.1.0
```

#### docker 驗證設定
目前有一些限制 , 需要手動通過驗證

- Claude code by Claude.ai account: Claude.ai 目前沒有 API token 可以自動化登入
  - `docker exec -it autodev-agent claude` 執行 claude code , 並根據指示進入交互模式進行驗證
  - 點擊驗證連結 , 進行驗證後將驗證碼貼回Claude Code

- Notion MCP: Notion MCP 目前沒有自動化登入
  - 在Claude Code 交互模式 `/mcp` , 選擇 `notion`
  - 選擇Authenticate , 進入驗證 , 交互會顯示驗證連結 , 複製驗證連結到瀏覽器進行驗證
  - 完成後會轉跳 , 會出現找不到網頁 , 將callback網址複製
  - 輸入以下指令 , 完成驗證
  ```
  docker exec -it autodev-agent curl -X GET "<call_back_url>"
  ```

#### Local development
```bash
# Clone repo
git clone github.com/itmrchow/autodev-agent
cd autodev-agent

# 複製環境變量範本
cp .env.example .env
# 編輯 .env 填入相關 tokens

# 啟動服務
go run main.go
```

### 健康檢查
服務啟動後可以透過以下端點檢查狀態：
```bash
curl http://localhost:8090/health
```

## Todo-list
- other model or Agent framework
- claude code 持久化
- update prompt to English
- other pm tool
- Story ticket 分析
- Agent開發通知(開始&結束)
- Agent處理狀態
- 同時處理上限設定