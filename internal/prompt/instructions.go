package prompt

import "strings"

const (
	instructionsStartTag = "<instructions>"
	instructionsEndTag   = "</instructions>"
)

func buildInstructions(r *BuilderRequest) string {
	var builder strings.Builder

	builder.WriteString(instructionsStartTag)
	builder.WriteString("\n")
	builder.WriteString(evidenceBlock)
	builder.WriteString("\n")
	builder.WriteString(verdictsBlock)
	builder.WriteString("\n")
	builder.WriteString(skillsBlock)
	builder.WriteString("\n")
	builder.WriteString(buildTools(r.Tools))
	builder.WriteString("\n")
	builder.WriteString(filesBlock)
	builder.WriteString("\n")
	builder.WriteString(workingBlock)
	builder.WriteString(instructionsEndTag)
	builder.WriteString("\n")

	return builder.String()
}
