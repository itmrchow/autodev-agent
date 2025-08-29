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
- `PM_TOOL_NOTION_TOKEN` Notion API token
- `PM_TOOL_JIRA_TOKEN` Jira API token (optional)
- `CODE_REPO_GITHUB_TOKEN` GitHub API token
- `PORT` 服務端口，預設 `8090`

#### GitHub 配置（用於 Claude Code 開發功能）
- `GITHUB_USER_NAME` GitHub 使用者名稱
- `GITHUB_USER_EMAIL` GitHub 使用者 email
- `GITHUB_PERSONAL_TOKEN` GitHub Personal Access Token

#### Claude Code 配置
- `CLAUDE_API_KEY` Claude API key（用於 AI 開發功能）

### 本地開發
```bash
# Clone repo
git clone <repo-url>
cd autodev-agent

# 複製環境變量範本
cp .env.example .env
# 編輯 .env 填入相關 tokens

# 啟動服務
go run main.go
```

### Docker 部署

#### 1. 使用預構建的 Docker Image
```bash
# 拉取最新映像
docker pull <registry>/autodev-agent:latest

# 運行容器（包含完整配置）
docker run -d \
  -p 8090:8090 \
  -e PM_TOOL_NOTION_TOKEN=your_notion_token \
  -e CODE_REPO_GITHUB_TOKEN=your_github_token \
  -e GITHUB_USER_NAME="Your Name" \
  -e GITHUB_USER_EMAIL="your.email@example.com" \
  -e GITHUB_PERSONAL_TOKEN=your_github_personal_token \
  -e CLAUDE_API_KEY=your_claude_api_key \
  -v $(pwd)/projects:/app/projects \
  --name autodev-agent \
  <registry>/autodev-agent:latest
```

#### 2. 本地構建 Docker Image
```bash
# 構建映像
./build.sh

# 或指定名稱和標籤
./build.sh -n myapp -t v1.0.0

# 運行容器（包含完整配置）
docker run -d \
  -p 8090:8090 \
  -e PM_TOOL_NOTION_TOKEN=your_notion_token \
  -e CODE_REPO_GITHUB_TOKEN=your_github_token \
  -e GITHUB_USER_NAME="Your Name" \
  -e GITHUB_USER_EMAIL="your.email@example.com" \
  -e GITHUB_PERSONAL_TOKEN=your_github_personal_token \
  -e CLAUDE_API_KEY=your_claude_api_key \
  -v $(pwd)/projects:/app/projects \
  --name autodev-agent \
  autodev-agent:latest
```

#### 3. 發布到 Docker Registry
```bash
# 發布到 Docker Hub
./publish.sh -r yourusername -t v1.0.0

# 發布到私有 registry
./publish.sh -r registry.com/yourusername -t v1.0.0
```

### GitHub Personal Access Token 設定

1. 前往 GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. 點擊 "Generate new token (classic)"
3. 選擇所需權限：
   - `repo` - 完整的 repository 權限
   - `workflow` - 如果需要操作 GitHub Actions
4. 複製產生的 token 並設定到 `GITHUB_PERSONAL_TOKEN` 環境變數

### Claude API Key 設定

1. 前往 [Anthropic Console](https://console.anthropic.com/)
2. 登入並前往 API Keys 頁面
3. 點擊 "Create Key" 建立新的 API key
4. 複製產生的 API key 並設定到 `CLAUDE_API_KEY` 環境變數

### 健康檢查
服務啟動後可以透過以下端點檢查狀態：
```bash
curl http://localhost:8090/health
```

## Todo-list
- Story ticket 分析
- Agent開發通知(開始&結束)
- Agent處理狀態
- 同時處理上限設定

