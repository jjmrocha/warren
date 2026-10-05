package prompt

import (
	"strings"

	"github.com/jjmrocha/ai-toolkit/mcp"
)

func buildInstructions(tools []mcp.Instruction) string {
	blocks := []string{
		evidenceBlock,
		verdictsBlock,
		skillsBlock,
		buildTools(tools),
		filesBlock,
		workingBlock,
	}

	return "<instructions>\n" + strings.Join(blocks, "\n") + "</instructions>\n"
}
