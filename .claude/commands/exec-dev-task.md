執行開發任務

# 任務內容輸入
- 根據對話內容執行任務, 並回傳changelog
- {status}: 執行任務要更新的ticket狀態

# git
分支命名規則
- 功能開發：feat/ticket_id/功能描述
- Bug 修復：bugfix/ticket_id/問題描述
- 熱修復：hotfix/緊急修復描述
- 重構：refactor/ticket_id/重構內容

## ticket_id
- notion: ticket_id = id
- jira: ticket_id = key

# task list
- task exec
  - 判斷ticket專案
    - 檢查projects資料夾與資料夾中的projects.md , 判斷是否存在於服務
      - 已存在的專案
        - 於該專案的資料夾內執行任務
      - 不存在專案
        - 於github尋找專案 , clone到專案資料夾
        - 設定並確認remote repo為github上尋找到的專案repo
        - 更新projects.md
        - 於該資料夾內執行任務 , 並確認git branch 於 main branch
  - 開發環境建立 , 並使用git worktree 建立開發資料夾
    - git 切換到 main branch
    - 拉取最新變更 git pull origin main
    - 建立開發環境
      - 建立開發分支 git branch [branch-name]
      - 建立開發資料夾 git worktree add worktree-[project-name] [branch-name]
  - 任務開發
    - 更新ticket狀態為 {status}
    - 進入開發資料夾 , 於開發資料夾內開發
    - 讀取專案資料夾內claude.md , 根據claude.md內的指令進行執行
    - 根據task內容進行功能開發
    - 可以跳過功能測試部分, 後續review時再進行測試
  - 任務完成
    - git push origin [branch-name]
    - git worktree remove worktree-[project-name]
- change log
  - 根據跟異動內容產生change log

# 備註
- 如果專案的claude.md有指示等待確認,則一率忽略等待指示直接執行,本command是自動化流程,不需要人工確認
