package usecase

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ExecClaudeCode struct {
	rootPath string
	logger   zerolog.Logger
}

func NewExecClaudeCode(logger zerolog.Logger) *ExecClaudeCode {
	wd, _ := os.Getwd()

	return &ExecClaudeCode{
		rootPath: wd,
		logger:   logger,
	}
}

func (c *ExecClaudeCode) ExecTask(ticketId string, pmTool string) error {

	var execFileName string

	if pmTool == "notion" {
		execFileName = "exec-notion-task.md"
	} else {
		return errors.New("pmTool not support")
	}

	c.logger.Info().Str("ticket_id", ticketId).Msg("開始執行 Claude 命令")

	commandCmd := exec.Command("cat", ".claude/commands/"+execFileName)
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
