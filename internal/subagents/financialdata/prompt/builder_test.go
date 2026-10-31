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

	t.Run("orders the blocks: tools, fetching, answer", func(t *testing.T) {
		// given
		request := Request{}
		// when
		result := Build(request)
		// then
		tools := strings.Index(result, "<tools>")
		fetching := strings.Index(result, "<how-to-fetch>")
		answer := strings.Index(result, "<answer>")
		assert.Positive(t, tools)
		assert.Less(t, tools, fetching)
		assert.Less(t, fetching, answer)
	})

	t.Run("never guesses a ticker", func(t *testing.T) {
		// given
		request := Request{}
		// when
		result := Build(request)
		// then
		assert.Contains(t, result, "never guess a ticker from a company name")
	})

	t.Run("copies figures as returned and says they are Yahoo's", func(t *testing.T) {
		// given
		request := Request{}
		// when
		result := Build(request)
		// then
		assert.Contains(t, result, "- Copy every figure exactly as the tool returns it, with its unit, currency\n  and period")
		assert.Contains(t, result, "- Say that the figures are Yahoo Finance's, not the company's own filing.")
	})

	t.Run("names no location, because the financial-data agent reaches no folder", func(t *testing.T) {
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
				{Name: "yfinance", Text: "Quotes may be delayed."},
				{Name: "archive", Text: "Read old data."},
			},
		}
		// when
		result := Build(request)
		// then
		expected := "<tools>\n" +
			"<tool-instructions name=\"archive\">\nRead old data.\n</tool-instructions>\n" +
			"<tool-instructions name=\"yfinance\">\nQuotes may be delayed.\n</tool-instructions>\n" +
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
