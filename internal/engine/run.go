package engine

import (
	"context"

	"github.com/jjmrocha/ai-chat/chat"
	"github.com/jjmrocha/ai-chat/ui"
	"github.com/jjmrocha/ai-toolkit/agent"
	"github.com/jjmrocha/ai-toolkit/llm"
	"github.com/jjmrocha/ai-toolkit/packs"
	"github.com/jjmrocha/ai-toolkit/tools"
	"github.com/jjmrocha/warren/internal/config"
	"github.com/jjmrocha/warren/internal/prompt"
)

func Run(ctx context.Context, cfg *config.Config) error {
	// Initialize the LLM
	llmClient, err := llm.New(cfg.LLMConfig())
	if err != nil {
		return err
	}

	// Initialize the skills collection
	skills, err := newSkillCollection(cfg)
	if err != nil {
		return err
	}

	// Initialize the toolbox
	toolBox := tools.NewToolBox()

	// Initialize the MCP manager
	mng := newMCPManager(toolBox, cfg)

	defer mng.Close()

	// Start the MCP servers the config boots
	startMCPs(ctx, mng, cfg)

	// Register tools
	filePack, err := packs.FileTools(toolBox, ".")
	if err != nil {
		return err
	}

	defer func() { _ = filePack.Close() }()

	webPack, err := packs.WebTools(ctx, toolBox)
	if err != nil {
		return err
	}

	defer func() { _ = webPack.Close() }()

	datePack, err := packs.DateTools(toolBox)
	if err != nil {
		return err
	}

	defer func() { _ = datePack.Close() }()

	// Initialize the agent
	ag, err := agent.New(agent.Config{}, llmClient)
	if err != nil {
		return err
	}

	defer ag.Close()

	// Initialize the chat
	chatAgent := chat.New("WARREN", ag,
		chat.WithDefaultCommands(),
		chat.WithMCP(mng),
		chat.WithSkills(skills),
	)

	// Set session
	ag.StartSession(agent.SessionConfig{
		Prompt:  prompt.Build(),
		Skills:  skills,
		ToolBox: toolBox,
	})

	return ui.Run(ctx, chatAgent)
}
