package controllers

import (
	"strconv"

	"github.com/TimothyYe/godns/pkg/lib"
	"github.com/gofiber/fiber/v2"
)

// GetLogs returns log entries.
func (c *Controller) GetLogs(ctx *fiber.Ctx) error {
	// Get query parameters
	limitStr := ctx.Query("limit", "100")
	level := ctx.Query("level", "")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	// Get log buffer
	logBuffer := lib.GetLogBuffer()
	entries := logBuffer.GetRecent(limit)

	// Filter by level if specified
	if level != "" {
		filteredEntries := make([]lib.LogEntry, 0)
		for _, entry := range entries {
			if entry.Level == level {
				filteredEntries = append(filteredEntries, entry)
			}
		}
		entries = filteredEntries
	}

	return ctx.JSON(fiber.Map{
		"logs":  entries,
		"total": len(entries),
	})
}

// ClearLogs clears all log entries.
func (c *Controller) ClearLogs(ctx *fiber.Ctx) error {
	logBuffer := lib.GetLogBuffer()
	logBuffer.Clear()

	return ctx.JSON(fiber.Map{
		"message": "Logs cleared successfully",
	})
}

// GetLogLevels returns available log levels.
func (c *Controller) GetLogLevels(ctx *fiber.Ctx) error {
	levels := []string{"panic", "fatal", "error", "warn", "info", "debug", "trace"}

	return ctx.JSON(fiber.Map{
		"levels": levels,
	})
}
