package ui

import (
	"strings"

	"github.com/fatih/color"
)

const (
	// Colors
	SuccessColor = color.FgGreen
	ErrorColor   = color.FgRed
	InfoColor    = color.FgCyan
	WarningColor = color.FgYellow
)

// Color functions
var (
	SuccessPrint = color.New(SuccessColor).SprintFunc()
	ErrorPrint   = color.New(ErrorColor).SprintFunc()
	InfoPrint    = color.New(InfoColor).SprintFunc()
	WarningPrint = color.New(WarningColor).SprintFunc()
)

func GetStatusColor(status string) func(...interface{}) string {
	switch strings.ToLower(status) {
	case "running":
		return color.New(color.FgGreen).SprintFunc()
	case "deploying":
		return color.New(color.FgYellow).SprintFunc()
	case "failed", "error":
		return color.New(color.FgRed).SprintFunc()
	case "pending":
		return color.New(color.FgYellow).SprintFunc()
	case "crash_loop":
		return color.New(color.FgRed, color.Bold).SprintFunc()
	default:
		return color.New(color.FgWhite).SprintFunc()
	}
}

func GetLogLevelColor(level string) func(...interface{}) string {
	switch strings.ToLower(level) {
	case "error":
		return color.New(color.FgRed).SprintFunc()
	case "warn":
		return color.New(color.FgYellow).SprintFunc()
	case "info":
		return color.New(color.FgCyan).SprintFunc()
	case "debug":
		return color.New(color.FgWhite).SprintFunc()
	default:
		return color.New(color.FgWhite).SprintFunc()
	}
}
