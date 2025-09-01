#!/bin/bash

# Git 初始化腳本

set -e

echo "🔧 Initializing Git configuration..."

# 檢查必要的環境變數
missing_vars=()

if [ -z "$GITHUB_USER_NAME" ]; then
    missing_vars+=("GITHUB_USER_NAME")
fi

if [ -z "$GITHUB_USER_EMAIL" ]; then
    missing_vars+=("GITHUB_USER_EMAIL")
fi

if [ -z "$GITHUB_PERSONAL_TOKEN" ]; then
    missing_vars+=("GITHUB_PERSONAL_TOKEN")
fi

# 如果有缺少的環境變數，顯示警告但繼續執行
if [ ${#missing_vars[@]} -gt 0 ]; then
    echo "⚠️  Missing GitHub environment variables: ${missing_vars[*]}"
    echo "   Git functionality will be limited without these variables:"
    echo "   - GITHUB_USER_NAME: Your GitHub username"
    echo "   - GITHUB_USER_EMAIL: Your GitHub email"
    echo "   - GITHUB_PERSONAL_TOKEN: Your GitHub personal access token"
    echo "   Skipping Git configuration..."
    return 0 2>/dev/null || exit 0
fi

echo "✅ All GitHub environment variables provided"

# 設定 Git 全域配置
echo "📝 Setting Git global configuration..."
git config --global user.name "$GITHUB_USER_NAME"
git config --global user.email "$GITHUB_USER_EMAIL"

# 設定 Git credential helper 使用 personal token
echo "🔑 Setting up Git credential helper..."

# 創建 credential helper 腳本
cat > /tmp/git-credential-helper.sh << 'EOF'
#!/bin/bash
echo username=$GITHUB_USER_NAME
echo password=$GITHUB_PERSONAL_TOKEN
EOF

chmod +x /tmp/git-credential-helper.sh

# 設定 Git 使用我們的 credential helper
git config --global credential.helper '!/tmp/git-credential-helper.sh'

# 設定一些實用的 Git 配置
git config --global init.defaultBranch main
git config --global pull.rebase false
git config --global core.autocrlf input
git config --global color.ui auto

# 驗證配置
echo "🔍 Verifying Git configuration..."
echo "   User Name: $(git config --global user.name)"
echo "   User Email: $(git config --global user.email)"
echo "   Default Branch: $(git config --global init.defaultBranch)"

# 測試 GitHub 連接（可選，不影響啟動）
echo "🌐 Testing GitHub connectivity..."
if git ls-remote https://github.com/octocat/Hello-World.git >/dev/null 2>&1; then
    echo "✅ GitHub connectivity test successful"
else
    echo "⚠️  GitHub connectivity test failed - please check your token"
fi

echo "✅ Git initialization completed!"