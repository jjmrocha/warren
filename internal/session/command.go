package session

import (
	"context"
	"fmt"
	"slices"

	"github.com/jjmrocha/ai-chat/command"
	"github.com/jjmrocha/ai-toolkit/agent"
	"github.com/jjmrocha/ai-toolkit/llm"
)

type Source interface {
	Messages() []llm.Message
	ModelInfo(ctx context.Context) *agent.ModelInfo
	SessionID() string
}

type exportCmd struct {
	src Source
}

func ExportCommand(src Source) command.Command {
	return exportCmd{src: src}
}

func (exportCmd) Name() string {
	return "export"
}

func (exportCmd) Help() string {
	return "Export session"
}

func (c exportCmd) Run(ctx command.Context, _ string) {
	msgs := c.src.Messages()
	if !slices.ContainsFunc(msgs, isTurn) {
		ctx.Print(command.Info, "No session to export.")
		return
	}

	if err := export(ctx.Context(), c.src, msgs); err != nil {
		ctx.Print(command.Error, fmt.Sprintf("export: %v", err))
		return
	}

	ctx.Print(command.Info, fmt.Sprintf("Session %s exported", c.src.SessionID()))
}

func isTurn(m llm.Message) bool {
	return m.Role() != llm.SystemRole
}
