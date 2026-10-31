package prompt

import (
	"strings"
	"testing"

	"github.com/jjmrocha/ai-toolkit/mcp"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	t.Run("starts with the role, followed by the instructions", func(t *testing.T) {
		// given
		request := Request{}
		// when
		result := Build(request)
		// then
		assert.True(t, strings.HasPrefix(result, "\n<role>\n"))
		assert.Less(t, strings.Index(result, "</role>"), strings.Index(result, "<instructions>"))
		assert.True(t, strings.HasSuffix(result, "</instructions>\n"))
	})

	t.Run("orders the blocks: tools, reading, answer", func(t *testing.T) {
		// given
		request := Request{}
		// when
		result := Build(request)
		// then
		tools := strings.Index(result, "<tools>")
		reading := strings.Index(result, "<how-to-read>")
		answer := strings.Index(result, "<answer>")
		assert.Positive(t, tools)
		assert.Less(t, tools, reading)
		assert.Less(t, reading, answer)
	})

	t.Run("copies figures as printed and cites their document, period and URL", func(t *testing.T) {
		// given
		request := Request{}
		// when
		result := Build(request)
		// then
		assert.Contains(t, result, "- Copy every figure exactly as the document prints it, with its unit, currency\n  and period. Never round, convert or compute one.")
		assert.Contains(t, result, "- Give every figure and quote its document, period and URL")
	})

	t.Run("names no location, because the web-reader reaches no folder", func(t *testing.T) {
		// given
		request := Request{}
		// when
		result := Build(request)
		// then
		assert.NotContains(t, result, "<locations>")
	})

	t.Run("renders each tool's instructions under its name, in name order", func(t *testing.T) {
		// given
		request := Request{
			ToolInstructions: []mcp.Instruction{
				{Name: "donsetch", Text: "Search the web."},
				{Name: "browser", Text: "Open a page."},
			},
		}
		// when
		result := Build(request)
		// then
		expected := "<tools>\n" +
			"<tool-instructions name=\"browser\">\nOpen a page.\n</tool-instructions>\n" +
			"<tool-instructions name=\"donsetch\">\nSearch the web.\n</tool-instructions>\n" +
			"</tools>\n"
		assert.Contains(t, result, expected)
	})

	t.Run("leaves no format verb unfilled", func(t *testing.T) {
		// given
		request := Request{}
		// when
		result := Build(request)
		// then
		assert.NotContains(t, result, "%!")
		assert.NotContains(t, result, "%s")
	})
}
