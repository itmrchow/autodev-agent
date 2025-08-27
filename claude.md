# autodev-agent

這個服務是一個 SPEC 驅動的 webhook 處理器，負責接收來自專案管理工具的事件，並觸發 AI Agent 進行任務

## tech stack
- go + gin : webhook
- pm tool
  - notion
- ai agent
  - claude

## claude 使用場景與流程
- webhook開發
- webhook呼叫執行任務
  - 會呼叫 commands , 以commands內的指令執行
  <!-- - TODO: 檢查是否有專案共通的rules.md -->


## Webhook呼叫執行任務開發策略

### 專案開發 Git Repo 使用原則
- **主 Repo (autodev-agent)**：僅用於專案管理和整體協調
- **實際開發**：需要從 GitHub 搜尋對應的專案 repo 進行開發

### 開發流程
1. 不要直接在 autodev-agent 中進行功能開發 , 只能在projects資料夾內進行 , 並且在git worktree內新增開發內容
2. 根據任務類型搜尋對應的 GitHub repo
3. Clone 對應 repo 到 projects/[專案名]/[模組名] 資料夾
4. 在對應的 repo 中建立分支和進行開發
5. 推送變更到對應的 repo，而非主 repo