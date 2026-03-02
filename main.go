package main

import (
	"log/slog"
	"kassemblycodegen/internal/command"
)

func main() {
	if err := command.Execute(); err != nil {
		slog.Error("command execution failed", "error", err)
	}
}
