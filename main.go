package main

import (
	"kassemblycodegen/internal/command"
	"log/slog"
	"os"
)

func setlog() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, opts))
	slog.SetDefault(logger)
}

func main() {
	setlog()
	if err := command.Execute(); err != nil {
		slog.Error("command execution failed", "error", err)
	}
}
