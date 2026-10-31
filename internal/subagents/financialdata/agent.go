package financialdata

import (
	"context"
	"fmt"

	"github.com/jjmrocha/ai-toolkit/agent"
	"github.com/jjmrocha/ai-toolkit/llm"
	"github.com/jjmrocha/ai-toolkit/mcp"
	"github.com/jjmrocha/ai-toolkit/packs"
	"github.com/jjmrocha/ai-toolkit/subagent"
	"github.com/jjmrocha/ai-toolkit/tools"
	"github.com/jjmrocha/warren/internal/subagents/financialdata/prompt"
)

const (
	name        = "financial-data"
	description = "Reads market and accounting data from Yahoo Finance — quotes, price history, financial statements, analyst estimates and holders — and returns the figures the task asks for, each with its period and currency. The task must give the ticker, with its exchange suffix outside the US (for example ASML.AS)."
)

var financeTools = packs.FinanceTools

func Register(ctx context.Context, manager *subagent.Manager, llmConfig llm.Config) (closeFn func(), err error) {
	// Initialize the model
	llmClient, err := llm.New(llmConfig)
	if err != nil {
		return nil, fmt.Errorf("model: %w", err)
	}

	// Register tools
	toolBox := tools.NewToolBox()

	financePack, err := financeTools(ctx, toolBox)
	if err != nil {
		return nil, fmt.Errorf("finance tools: %w", err)
	}

	closeFn = func() { _ = financePack.Close() }

	defer func() {
		if err != nil {
			closeFn()
		}
	}()

	// Initialize the agent
	financialDataAgent, err := agent.New(agent.Config{}, llmClient)
	if err != nil {
		return nil, fmt.Errorf("agent: %w", err)
	}

	// Build prompt
	var toolInstructions []mcp.Instruction

	if instruction := financePack.Instructions(ctx); instruction != nil {
		toolInstructions = append(toolInstructions, *instruction)
	}

	// Register the financial-data agent
	financialDataConfig := subagent.Config{
		Name:        name,
		Description: description,
		Prompt:      prompt.Build(prompt.Request{ToolInstructions: toolInstructions}),
	}

	if err = manager.Register(financialDataConfig, financialDataAgent, toolBox, nil); err != nil {
		return nil, err
	}

	return closeFn, nil
}
