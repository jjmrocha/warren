package prompt

import "github.com/jjmrocha/ai-toolkit/mcp"

func Build(tools []mcp.Instruction) string {
	return rolePrompt + buildInstructions(tools)
}
