package prompt

import (
	"strings"

	"github.com/jjmrocha/ai-toolkit/mcp"
)

type BuilderRequest struct {
	Tools []mcp.Instruction
}

func Build(r *BuilderRequest) string {
	var builder strings.Builder

	builder.WriteString(rolePrompt)
	builder.WriteString(buildInstructions(r))

	return builder.String()
}
