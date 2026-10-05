package engine

import (
	"context"
	"fmt"

	"github.com/jjmrocha/ai-chat/chat"
	"github.com/jjmrocha/ai-chat/ui"
	"github.com/jjmrocha/ai-toolkit/agent"
	"github.com/jjmrocha/ai-toolkit/llm"
	"github.com/jjmrocha/ai-toolkit/packs"
	"github.com/jjmrocha/ai-toolkit/tools"
	"github.com/jjmrocha/warren/internal/config"
	"github.com/jjmrocha/warren/internal/prompt"
	"github.com/jjmrocha/warren/internal/session"
)

func Run(ctx context.Context, cfg *config.Config, sessionID string) error {
	// Initialize the LLM
	llmClient, err := llm.New(cfg.LLMConfig())
	if err != nil {
		return fmt.Errorf("model: %w", err)
	}

	// Initialize the skills collection
	skills, err := newSkillCollection(cfg)
	if err != nil {
		return fmt.Errorf("skills: %w", err)
	}

	// Restore session
	var history []llm.Message

	if sessionID != "" {
		history, err = session.Load(sessionID)
		if err != nil {
			return fmt.Errorf("resume: %w", err)
		}
	}

	// Initialize the toolbox
	toolBox := tools.NewToolBox()
	toolBox.SetInterceptor(guardFiles)

	// Initialize the MCP manager
	mng := newMCPManager(toolBox, cfg)

	defer mng.Close()

	// Start the MCP servers the config boots
	startMCPs(ctx, mng, cfg)

	// Register tools
	var toolPacks packSet

	defer toolPacks.close()

	if err = toolPacks.add(packs.FileTools(toolBox, ".")); err != nil {
		return fmt.Errorf("file tools: %w", err)
	}

	if err = toolPacks.add(packs.WebTools(ctx, toolBox)); err != nil {
		return fmt.Errorf("web tools: %w", err)
	}

	if err = toolPacks.add(packs.DateTools(toolBox)); err != nil {
		return fmt.Errorf("date tools: %w", err)
	}

	// Initialize the agent
	ag, err := agent.New(agent.Config{}, llmClient)
	if err != nil {
		return fmt.Errorf("agent: %w", err)
	}

	defer ag.Close()

	// Initialize the chat
	chatAgent := chat.New("WARREN", ag, buildCommands(mng, skills, ag)...)

	// Build prompt
	systemPrompt := prompt.Build(toolInstructions(ctx, toolPacks, mng))

	// Set session
	ag.StartSession(agent.SessionConfig{
		Prompt:   systemPrompt,
		Skills:   skills,
		ToolBox:  toolBox,
		Messages: history,
		ID:       sessionID,
	})

	return ui.Run(ctx, chatAgent)
}
