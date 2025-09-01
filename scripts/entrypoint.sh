#!/bin/bash

# Docker 容器主啟動腳本

set -e

echo "🚀 Starting autodev-agent container..."
echo "========================================"

# 檢查並創建必要的目錄
echo "📁 Setting up directories..."
mkdir -p /tmp

# 檢查核心環境變量
echo "🔧 Checking environment variables..."
echo "   PORT: $PORT"
echo "   GIN_MODE: $GIN_MODE"

# 檢查 .claude 配置
if [ -d "/app/.claude" ]; then
    echo "✅ Claude Code configuration found"
    command_count=$(find /app/.claude/commands -name "*.md" 2>/dev/null | wc -l)
    echo "   Available commands: $command_count"
else
    echo "⚠️  No .claude directory found"
fi

# 執行 Git 初始化
echo ""
echo "🔧 Git Configuration"
echo "===================="
if [ -f "/app/init-git.sh" ]; then
    source /app/init-git.sh
else
    echo "⚠️  Git initialization script not found"
fi

# 執行 Claude Code 初始化
echo ""
echo "🤖 Claude Code Configuration"
echo "=========================="
if [ -n "$CLAUDE_API_KEY" ]; then
    echo "✅ Claude API key provided"
    
    # 設置 Claude API key
    export ANTHROPIC_API_KEY="$CLAUDE_API_KEY"
    
    # 驗證 Claude Code 是否可用
    if command -v claude >/dev/null 2>&1; then
        echo "✅ Claude Code CLI available"
        # 測試 Claude Code 連接
        if claude --version >/dev/null 2>&1; then
            echo "   Claude Code version: $(claude --version 2>/dev/null || echo 'Unable to get version')"
        else
            echo "⚠️  Claude Code CLI installed but may have issues"
        fi
    else
        echo "⚠️  Claude Code CLI not found"
    fi
else
    echo "⚠️  No Claude API key provided"
    echo "   Claude Code functionality will not be available"
    echo "   Set CLAUDE_API_KEY environment variable to enable Claude Code"
fi

# 顯示系統信息
echo ""
echo "📊 System Information"
echo "===================="
echo "   Working Directory: $(pwd)"
echo "   User: $(whoami)"
echo "   Available disk space: $(df -h /app 2>/dev/null | tail -1 | awk '{print $4}' || echo 'N/A')"
echo "   Git version: $(git --version 2>/dev/null || echo 'Git not available')"
echo "   Node.js version: $(node --version 2>/dev/null || echo 'Node.js not available')"
echo "   Claude Code: $(command -v claude >/dev/null 2>&1 && echo 'Available' || echo 'Not available')"

echo ""
echo "✅ Initialization completed!"
echo "🎯 Starting webhook service on port ${PORT}..."
echo "========================================"

# 啟動主服務
exec "$@"