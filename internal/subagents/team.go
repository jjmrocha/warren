package subagents

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/jjmrocha/ai-chat/chat"
	"github.com/jjmrocha/ai-toolkit/packs"
	"github.com/jjmrocha/ai-toolkit/subagent"
	"github.com/jjmrocha/ai-toolkit/tools"
	"github.com/jjmrocha/warren/internal/config"
	"github.com/jjmrocha/warren/internal/subagents/financialdata"
	"github.com/jjmrocha/warren/internal/subagents/webreader"
)

var (
	registerWebReader     = webreader.Register
	registerFinancialData = financialdata.Register
)

type Team struct {
	manager      *subagent.Manager
	subagentTool packs.ToolPack
	closers      []func()
	once         sync.Once
}

type Progress interface {
	SetText(text string)
}

func Start(ctx context.Context, toolBox *tools.ToolBox, cfg *config.Config, progress Progress) (*Team, error) {
	team := &Team{}

	if !cfg.Subagents.WebReader.Enabled && !cfg.Subagents.FinancialData.Enabled {
		return team, nil
	}

	manager := subagent.NewManager()

	if cfg.Subagents.WebReader.Enabled {
		progress.SetText("Starting subagent web-reader…")

		closeWebReader, err := registerWebReader(ctx, manager, cfg.SubagentLLMConfig(cfg.Subagents.WebReader))
		if err != nil {
			return nil, fmt.Errorf("web-reader: %w", err)
		}

		team.closers = append(team.closers, closeWebReader)
	}

	if cfg.Subagents.FinancialData.Enabled {
		progress.SetText("Starting subagent financial-data…")

		closeFinancialData, err := registerFinancialData(ctx, manager, cfg.SubagentLLMConfig(cfg.Subagents.FinancialData))
		if err != nil {
			team.Close()
			return nil, fmt.Errorf("financial-data: %w", err)
		}

		team.closers = append(team.closers, closeFinancialData)
	}

	subagentTool, err := packs.SubAgentTools(toolBox, manager)
	if err != nil {
		team.Close()
		return nil, err
	}

	team.manager = manager
	team.subagentTool = subagentTool

	return team, nil
}

func (t *Team) ChatOption() chat.Option {
	return chat.WithSubagents(t.manager)
}

func (t *Team) Close() {
	t.once.Do(func() {
		if t.subagentTool != nil {
			_ = t.subagentTool.Close()
		}

		for _, closeFn := range slices.Backward(t.closers) {
			closeFn()
		}
	})
}
