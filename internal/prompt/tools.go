package prompt

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jjmrocha/ai-toolkit/mcp"
)

func buildTools(instructions []mcp.Instruction) string {
	var builder strings.Builder

	builder.WriteString("<tools>\n")

	sorted := slices.SortedFunc(slices.Values(instructions), func(a, b mcp.Instruction) int {
		return strings.Compare(a.Name, b.Name)
	})

	for _, instruction := range sorted {
		fmt.Fprintf(&builder, "<tool-instructions name=%q>\n%s\n</tool-instructions>\n", instruction.Name, instruction.Text)
	}

	builder.WriteString("</tools>\n")

	return builder.String()
}
