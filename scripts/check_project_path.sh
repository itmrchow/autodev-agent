#!/bin/bash

# 檢查 .env 檔案是否存在
if [ ! -f ".env" ]; then
  echo "錯誤：找不到 .env 檔案。"
  exit 1
fi

# 從 .env 載入環境變數
export $(grep -v '^#' .env | xargs)

# 檢查 PROJECT_PATH 是否已設定
if [ -z "$PROJECT_PATH" ]; then
  echo "錯誤：請在 .env 檔案中設定 PROJECT_PATH。"
  exit 1
fi

echo "目標路徑: $PROJECT_PATH"

# 檢查路徑是否存在，若不存在則建立
if [ -d "$PROJECT_PATH" ]; then
  echo "✅ 資料夾已存在。"
else
  echo "資料夾不存在，嘗試建立..."
  mkdir -p "$PROJECT_PATH"
  if [ $? -ne 0 ]; then
    echo "❌ 錯誤：無法建立資料夾 $PROJECT_PATH。"
    exit 1
  fi
  echo "✅ 資料夾已成功建立。"
fi

# 檢查 projects.md 檔案是否存在，若不存在則建立
PROJECT_MD_PATH="$PROJECT_PATH/projects.md"
if [ -f "$PROJECT_MD_PATH" ]; then
  echo "✅ 找到 projects.md 檔案。"
else
  echo "projects.md 不存在，嘗試建立..."
  touch "$PROJECT_MD_PATH"
  if [ $? -ne 0 ]; then
    echo "❌ 錯誤：無法建立檔案 $PROJECT_MD_PATH。"
    exit 1
  fi
  echo "✅ projects.md 已成功建立。"
fi

echo "🎉 檢查成功！"
exit 0
