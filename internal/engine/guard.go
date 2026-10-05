package engine

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jjmrocha/ai-toolkit/llm"
	"github.com/jjmrocha/warren/internal/config"
	"github.com/jjmrocha/warren/internal/session"
)

var fileWriters = []string{"file_write", "file_edit", "file_delete"}

var protectedPaths = []string{config.FileName, session.DirName}

func guardFiles(_ context.Context, call llm.ToolCall) error {
	if !slices.Contains(fileWriters, call.Name) {
		return nil
	}

	path, _ := call.Arguments["path"].(string)
	top, _, _ := strings.Cut(filepath.ToSlash(filepath.Clean(path)), "/")

	if slices.ContainsFunc(protectedPaths, func(protected string) bool { return strings.EqualFold(top, protected) }) {
		return fmt.Errorf("%w: %s", ErrProtectedPath, path)
	}

	return nil
}
