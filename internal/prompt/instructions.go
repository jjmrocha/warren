package prompt

import "strings"

const (
	instructionsStartTag = "<instructions>"
	instructionsEndTag   = "</instructions>"
)

func buildInstructions(r Request) string {
	var builder strings.Builder

	builder.WriteString("\n")
	builder.WriteString(instructionsStartTag)
	builder.WriteString("\n")

	builder.WriteString(buildEvidence(r.WithWebReader))

	if r.WithWebReader {
		builder.WriteString(buildWebReader())
	}

	builder.WriteString(buildVerdicts())
	builder.WriteString(buildSkills())
	builder.WriteString(buildTools(r.ToolInstructions))
	builder.WriteString(buildFiles())
	builder.WriteString(buildWorkingWithUser())

	builder.WriteString(instructionsEndTag)
	builder.WriteString("\n")

	return builder.String()
}
