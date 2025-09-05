package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"webhook/internal/infrastructure/config"
)

type ExecClaudeCode struct {
	rootPath   string
	logger     zerolog.Logger
	jiraConfig *config.JiraConfig

	// commands path
	execNotionTaskPath   CommandPath
	readJiraTicketPath   CommandPath
	execDevTaskPath      CommandPath
	updateJiraTicketPath CommandPath
}

type CommandPath string

const (
	CommandPathExecNotionTask   CommandPath = ".claude/commands/exec-notion-task.md"
	CommandPathReadJiraTicket   CommandPath = ".claude/commands/read-jira-ticket.md"
	CommandPathExecDevTask      CommandPath = ".claude/commands/exec-dev-task.md"
	CommandPathUpdateJiraTicket CommandPath = ".claude/commands/update-jira-ticket.md"
)

func NewExecClaudeCode(logger zerolog.Logger, jiraConfig *config.JiraConfig) *ExecClaudeCode {
	wd, _ := os.Getwd()

	return &ExecClaudeCode{
		rootPath:   wd,
		logger:     logger,
		jiraConfig: jiraConfig,
		// commands path
		execNotionTaskPath:   CommandPathExecNotionTask,
		readJiraTicketPath:   CommandPathReadJiraTicket,
		execDevTaskPath:      CommandPathExecDevTask,
		updateJiraTicketPath: CommandPathUpdateJiraTicket,
	}
}

// TODO: rename to exec notion task , and use private method
func (c *ExecClaudeCode) ExecTask(ticketId string, pmTool string) error {

	if pmTool != "notion" {
		return errors.New("pmTool not support")
	}

	c.logger.Info().Str("ticket_id", ticketId).Msg("開始執行 Claude 命令")

	commandCmd := exec.Command("cat", string(c.execNotionTaskPath))
	commandCmd.Dir = c.rootPath

	claudeCmd := exec.Command(
		"claude",
		"--dangerously-skip-permissions",
		"-p", fmt.Sprintf("Execute these instructions with ticket_id='%s'", ticketId))
	claudeCmd.Dir = c.rootPath

	// 設定即時輸出到控制台
	claudeCmd.Stdout = os.Stdout
	claudeCmd.Stderr = os.Stderr

	claudeCmd.Stdin, _ = commandCmd.StdoutPipe()

	c.logger.Info().Str("ticket_id", ticketId).Msg("啟動 cat 命令")
	err := commandCmd.Start()
	if err != nil {
		return fmt.Errorf("啟動 cat 指令失敗: %w", err)
	}

	c.logger.Info().Str("ticket_id", ticketId).Msg("啟動 Claude 命令")
	err = claudeCmd.Start()
	if err != nil {
		return fmt.Errorf("啟動 claude 指令失敗: %w", err)
	}

	// 等待命令完成
	if err := commandCmd.Wait(); err != nil {
		log.Error().Err(err).Str("ticket_id", ticketId).Msg("cat 命令執行失敗")
	}

	if err := claudeCmd.Wait(); err != nil {
		return fmt.Errorf("claude 命令執行失敗: %w", err)
	}
	log.Info().Str("ticket_id", ticketId).Msg("Claude 命令執行完成")

	return nil
}

func (c *ExecClaudeCode) ExecJiraTask(ticketId string) error {

	// create session
	var sessionId string
	sessionId, err := c.exceCreateSession()
	if err != nil {
		c.logger.Err(err).Msg("Claude 建立 session 失敗")
		return err
	}

	// read ticket
	sessionId, err = c.exceReadJiraTicket(sessionId, ticketId)
	if err != nil {
		c.logger.Err(err).Msg("Claude 讀取 jira ticket 失敗")
		return err
	}

	// exec ticket -> ticket change log
	sessionId, err = c.execTask(sessionId, ticketId)
	if err != nil {
		c.logger.Err(err).Msg("Claude 執行 jira ticket 失敗")
		return err
	}

	// update ticket
	_, err = c.execUpdateJiraTicket(sessionId, ticketId, "") // 對話包含message
	if err != nil {
		c.logger.Err(err).Msg("Claude 更新 jira ticket 失敗")
		return err
	}

	return nil
}

// Claude Commands

// Create session , return session id
func (c *ExecClaudeCode) exceCreateSession() (string, error) {
	// Create session
	createSessionCmd := exec.Command(
		"claude",
		"-p", "In this conversation, I will provide a series of tasks, and you will be the executor, performing them based on the task content. If you receive this prompt, please respond with 'ok'.",
		"--output-format", "json",
	)
	createSessionCmd.Dir = c.rootPath

	output, err := createSessionCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("執行 claude 指令失敗: %w", err)
	}

	// 解析 JSON
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return "", fmt.Errorf("解析 JSON 失敗: %w", err)
	}

	sId := result["session_id"].(string)

	return sId, nil
}

// read jira ticket
func (c *ExecClaudeCode) exceReadJiraTicket(sessionId string, ticketId string) (newSessionId string, err error) {

	if sessionId == "" {
		return "", errors.New("sessionId is empty")
	}
	if ticketId == "" {
		return "", errors.New("ticketId is empty")
	}

	c.logger.Info().
		Str("ticket_id", ticketId).
		Str("session_id", sessionId).
		Msg("Claude 開始執行命令 - read-jira-ticket")

	readTicketCmd := exec.Command("cat", string(c.readJiraTicketPath))
	readTicketCmd.Dir = c.rootPath

	// pipe
	readCommandOutputPipe, err := readTicketCmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("command pipe failed: %w", err)
	}

	// claude command
	claudeCmd := exec.Command(
		"claude",
		"-r", sessionId,
		"--dangerously-skip-permissions",
		"--output-format", "json",
		"-p", fmt.Sprintf("Execute these instructions with ticket_id='%s'", ticketId),
	)
	claudeCmd.Dir = c.rootPath
	claudeCmd.Stdin = readCommandOutputPipe

	// start
	if err := readTicketCmd.Start(); err != nil {
		return "", fmt.Errorf("啟動 read 指令失敗: %w", err)
	}

	output, err := claudeCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("執行 claude 指令失敗: %w", err)
	}

	// wait
	if err := readTicketCmd.Wait(); err != nil {
		c.logger.Error().Err(err).Str("ticket_id", ticketId).Str("session_id", sessionId).Msg("cat 命令執行失敗")
		return "", fmt.Errorf("read 命令執行失敗: %w", err)
	}

	// 解析 JSON
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return "", fmt.Errorf("解析 JSON 失敗: %w", err)
	}

	newSessionId = result["session_id"].(string)

	log.Info().Str("ticket_id", ticketId).Str("session_id", newSessionId).Any("output", result["result"]).Msg("Claude 命令執行完成 - read-jira-ticket")

	return
}

// exec development task
func (c *ExecClaudeCode) execTask(sessionId string, ticketId string) (newSessionId string, err error) {

	if sessionId == "" {
		return "", errors.New("sessionId is empty")
	}

	c.logger.Info().
		Str("session_id", sessionId).
		Str("ticket_id", ticketId).
		Msg("Claude 開始執行命令 - exec-dev-task")

	execTaskCmd := exec.Command("cat", string(c.execDevTaskPath))
	execTaskCmd.Dir = c.rootPath

	// pipe
	execCommandOutputPipe, err := execTaskCmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("command pipe failed: %w", err)
	}

	prompt := fmt.Sprintf("Execute these instructions based on previous conversation and with status='%s'", c.jiraConfig.ExecStatus)

	// claude command
	claudeCmd := exec.Command(
		"claude",
		"-r", sessionId,
		"--dangerously-skip-permissions",
		"-p", prompt,
		"--output-format", "json",
	)
	claudeCmd.Dir = c.rootPath
	claudeCmd.Stdin = execCommandOutputPipe

	// start
	if err := execTaskCmd.Start(); err != nil {
		return "", fmt.Errorf("啟動 exec 指令失敗: %w", err)
	}

	output, err := claudeCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("執行 claude 指令失敗: %w", err)
	}

	// wait
	if err := execTaskCmd.Wait(); err != nil {
		c.logger.Error().Err(err).Str("session_id", sessionId).Msg("cat 命令執行失敗")
		return "", fmt.Errorf("exec 命令執行失敗: %w", err)
	}

	// 解析 JSON
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return "", fmt.Errorf("解析 JSON 失敗: %w", err)
	}

	newSessionId = result["session_id"].(string)

	log.Info().
		Str("session_id", sessionId).
		Str("ticket_id", ticketId).
		Any("output", string(output)).Msg("Claude 命令執行完成 - exec-dev-task")

	return
}

// update jira ticket
func (c *ExecClaudeCode) execUpdateJiraTicket(sessionId string, ticketId string, message string) (newSessionId string, err error) {

	c.logger.Info().
		Str("ticket_id", ticketId).
		Str("session_id", sessionId).
		Str("message", message).
		Msg("Claude 開始執行命令 - update-jira-ticket")

	updateTicketCmd := exec.Command("cat", string(c.updateJiraTicketPath))
	updateTicketCmd.Dir = c.rootPath

	// pipe
	updateCommandOutputPipe, err := updateTicketCmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("command pipe failed: %w", err)
	}

	// claude command with parameters
	prompt := fmt.Sprintf("Execute these instructions with ticket='%s', status='%s', message='%s'",
		ticketId, c.jiraConfig.ReviewStatus, message)

	claudeCmd := exec.Command(
		"claude",
		"-r", sessionId,
		"--dangerously-skip-permissions",
		"--output-format", "json",
		"-p", prompt,
	)
	claudeCmd.Dir = c.rootPath
	claudeCmd.Stdin = updateCommandOutputPipe

	// start
	if err := updateTicketCmd.Start(); err != nil {
		return "", fmt.Errorf("啟動 update 指令失敗: %w", err)
	}

	output, err := claudeCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("執行 claude 指令失敗: %w", err)
	}

	// wait
	if err := updateTicketCmd.Wait(); err != nil {
		c.logger.Error().Err(err).Str("ticket_id", ticketId).Str("session_id", sessionId).Msg("cat 命令執行失敗")
		return "", fmt.Errorf("update 命令執行失敗: %w", err)
	}

	// 解析 JSON
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return "", fmt.Errorf("解析 JSON 失敗: %w", err)
	}

	newSessionId = result["session_id"].(string)

	log.Info().Str("ticket_id", ticketId).Str("session_id", newSessionId).Any("output", result["result"]).Msg("Claude 命令執行完成 - update-jira-ticket")

	return
}
