# Webhook Build stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o autodev-agent ./main.go

#-----#

# Final stage
FROM alpine:latest

ARG CLAUDE_CODE_VERSION=latest

# 安裝必需的套件並建立用戶
RUN apk --no-cache add ca-certificates tzdata bash git nodejs npm && \
    addgroup -g 1001 appuser && \
    adduser -D -s /bin/sh -u 1001 -G appuser appuser && \
    rm -rf /var/cache/apk/*

# 設定工作目錄
WORKDIR /app

# 從 builder stage 複製二進制檔案
COPY --from=builder /app/autodev-agent .

# 安裝 Claude Code 並清理
RUN npm install -g @anthropic-ai/claude-code@${CLAUDE_CODE_VERSION} && \
    npm cache clean --force && \
    rm -rf /tmp/* /var/cache/apk/* && \
    apk del npm

COPY claude.md .mcp.json ./
COPY .claude/ .claude/
COPY scripts/ ./scripts/

# 修正 shell 腳本可能存在的 Windows 換行符問題 (CRLF -> LF)
RUN sed -i 's/\r$//' /app/scripts/*.sh && \ 
    chmod +x /app/scripts/*.sh && \
    chown -R appuser:appuser /app
  
# 切換到非 root 使用者
USER appuser

VOLUME [ "/app/projects" ]

# 設定環境變量
ENV GIN_MODE=release
ENV PORT=8090

# GitHub 相關環境變量（使用者需要在 docker run 時提供）
ENV GITHUB_USER_NAME=""
ENV GITHUB_USER_EMAIL=""
ENV GITHUB_PERSONAL_TOKEN=""

# PM Tool 相關環境變量
ENV PM_TOOL_NOTION_TOKEN=""
ENV PM_TOOL_JIRA_TOKEN=""

# Claude Code 相關環境變量
ENV CLAUDE_API_KEY=""

# 暴露端口
EXPOSE 8090

# 健康檢查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD nc -z localhost 8090 || exit 1

# 設定啟動腳本
ENTRYPOINT ["./scripts/entrypoint.sh"]

# 預設命令
CMD ["./autodev-agent"]