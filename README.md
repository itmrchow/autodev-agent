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

### Todo
- Story ticket 分析
- Agent開發通知(開始&結束)
- Agent處理狀態
- 同時處理上限設定

## tech stack

## 部署與設定

### Prerequisites

### Environment Variables
- PROJECT_PATH 專案資料夾路徑 , 預設./projects

### clone repo , 於本地啟動並使用claude code 相關設定

### docker image



