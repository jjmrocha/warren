package engine

import (
	"github.com/jjmrocha/ai-chat/chat"
	"github.com/jjmrocha/ai-toolkit/agent"
	"github.com/jjmrocha/ai-toolkit/mcp"
	"github.com/jjmrocha/ai-toolkit/skills"
	"github.com/jjmrocha/go-algo/fn"
	"github.com/jjmrocha/warren/internal/session"
)

func buildCommands(mng *mcp.Manager, skillCollection *skills.Collection, ag *agent.Agent) []chat.Option {
	options := fn.Map(coreSkills, func(skill coreSkill) chat.Option {
		return chat.WithSkillCommand(skill.name, skill.help)
	})

	return append(options,
		chat.WithDefaultCommands(),
		chat.WithMCP(mng),
		chat.WithSkills(skillCollection),
		chat.WithCommand(session.ExportCommand(ag)),
	)
}
