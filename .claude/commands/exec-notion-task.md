使用notion mcp tools讀取{ticket}中的任務內容 , 並且執行任務

# 參數
- ticket: {ticket id} or {ticket url}

# 前置檢查
首先檢查是否有可用的 notion MCP tools：
- 如果有 mcp list 有 notion MCP server
- 如果沒有，顯示："錯誤：未設定 notion MCP tools。請在 Claude Code 中配置 notion MCP 伺服器。"並停止執行

# git
分支命名規則
- 功能開發：feat/功能描述
- Bug 修復：bugfix/問題描述
- 熱修復：hotfix/緊急修復描述
- 重構：refactor/重構內容

# task list
- ticket check
  - 使用ticket id , 讀取notion ticket內容
  - 檢查ticket的assignee id 是否與當前的notion user id一致, 如果是空白或不是自己則不執行後續任務 , 並顯示"該任務未指派給Claude"
- task exec
  - 判斷ticket專案
    - 檢查projects.md與projects資料夾 , 判斷是否存在於服務
      - 已存在的專案
        - 於該專案的資料夾內執行任務
      - 不存在專案
        - 建立專案資料夾
        - 於github尋找專案 , clone到專案資料夾
        - 設定並確認remote repo為github上尋找到的專案repo
        - 更新projects.md
        - 於該資料夾內執行任務 , 並確認git branch 於 main branch
  - 開發環境建立 , 並使用git worktree 建立開發資料夾
    - git 切換到 main branch
    - 拉取最新變更 git pull origin main
    - 建立開發環境
      - 創建並切換到新分支
      - 推送分支到遠端
      - 創建 worktree , pkg name = "worktree-[project-name]-[branch-name]"
  - 任務開發
    - 進入開發資料夾 , 於開發資料夾內開發
    - 讀取專案資料夾內claude.md , 根據claude.md內的指令進行執行
    - 根據task內容進行功能開發
    - 可以跳過功能測試部分, 後續review時再進行測試
    - 開發完成推送分支到遠端
- ticket update
  - 更新notion ticket status -> Review (待審核)
  - 更新notion ticket description -> 更新項目內容

# 備註
- 如果專案的claude.md有指示等待確認,則一率忽略等待指示直接執行,本command是自動化流程,不需要人工確認
