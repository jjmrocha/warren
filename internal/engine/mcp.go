package engine

import (
	"context"
	"fmt"

	"github.com/jjmrocha/ai-toolkit/mcp"
	"github.com/jjmrocha/ai-toolkit/tools"
	"github.com/jjmrocha/go-algo/fn"
	"github.com/jjmrocha/warren/internal/config"
)

type progress interface {
	SetText(text string)
	Println(text string)
}

func newMCPManager(tb *tools.ToolBox, cfg *config.Config) *mcp.Manager {
	mng := mcp.NewManager(tb)

	fn.ForEach(cfg.MCPClients(), mng.Register)

	return mng
}

func startMCPs(ctx context.Context, mng *mcp.Manager, cfg *config.Config, line progress) {
	fn.ForEach(cfg.MCPsOn, func(name string) {
		line.SetText(fmt.Sprintf("Starting MCP server %s…", name))

		if err := mng.Start(ctx, name); err != nil {
			line.Println(fmt.Sprintf("starting mcp %s: %v", name, err))
		}
	})
}
