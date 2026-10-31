package webreader

import (
	"context"
	"fmt"

	"github.com/jjmrocha/ai-toolkit/agent"
	"github.com/jjmrocha/ai-toolkit/llm"
	"github.com/jjmrocha/ai-toolkit/mcp"
	"github.com/jjmrocha/ai-toolkit/packs"
	"github.com/jjmrocha/ai-toolkit/subagent"
	"github.com/jjmrocha/ai-toolkit/tools"
	"github.com/jjmrocha/warren/internal/subagents/webreader/prompt"
)

const (
	name        = "web-reader"
	description = "Searches the web and reads filings, reports, transcripts and other pages, and returns what the task asks for: figures and quotes with their document, period and URL, or a summary focused on the question."
)

var webTools = packs.WebTools

func Register(ctx context.Context, manager *subagent.Manager, llmConfig llm.Config) (closeFn func(), err error) {
	// Initialize the model
	llmClient, err := llm.New(llmConfig)
	if err != nil {
		return nil, fmt.Errorf("model: %w", err)
	}

	// Register tools
	toolBox := tools.NewToolBox()

	webPack, err := webTools(ctx, toolBox)
	if err != nil {
		return nil, fmt.Errorf("web tools: %w", err)
	}

	closeFn = func() { _ = webPack.Close() }

	defer func() {
		if err != nil {
			closeFn()
		}
	}()

	// Initialize the agent
	webReaderAgent, err := agent.New(agent.Config{}, llmClient)
	if err != nil {
		return nil, fmt.Errorf("agent: %w", err)
	}

	// Build prompt
	var toolInstructions []mcp.Instruction

	if instruction := webPack.Instructions(ctx); instruction != nil {
		toolInstructions = append(toolInstructions, *instruction)
	}

	// Register the web-reader
	webReaderConfig := subagent.Config{
		Name:        name,
		Description: description,
		Prompt:      prompt.Build(prompt.Request{ToolInstructions: toolInstructions}),
	}

	if err = manager.Register(webReaderConfig, webReaderAgent, toolBox, nil); err != nil {
		return nil, err
	}

	return closeFn, nil
}
