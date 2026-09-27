package prompt

import (
	"strings"
	"testing"

	"github.com/jjmrocha/ai-toolkit/mcp"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	t.Run("routes research ahead of valuation", func(t *testing.T) {
		// when
		result := Build(&BuilderRequest{})
		// then
		assert.Less(t, strings.Index(result, "company-research"), strings.Index(result, "buffett-valuation"))
	})

	t.Run("names the format and the shape of a report it is asked to write", func(t *testing.T) {
		// when
		result := Build(&BuilderRequest{})
		// then
		assert.Contains(t, result, "Markdown")
		assert.Contains(t, result, "subfolder")
	})

	t.Run("sends the model to the date tool instead of its own sense of the date", func(t *testing.T) {
		// when
		result := Build(&BuilderRequest{})
		// then
		assert.Contains(t, result, "current_date")
	})

	t.Run("tells the model to hand a skill's script an absolute path", func(t *testing.T) {
		// when
		result := Build(&BuilderRequest{})
		// then
		assert.Contains(t, result, "file_workdir")
	})

	t.Run("tells the model to remove a working file it no longer needs", func(t *testing.T) {
		// when
		result := Build(&BuilderRequest{})
		// then
		assert.Contains(t, result, "file_delete")
	})

	t.Run("names the tool that runs a skill's scripts", func(t *testing.T) {
		// when
		result := Build(&BuilderRequest{})
		// then
		assert.Contains(t, result, "skill_execute_file")
	})

	t.Run("never makes a missing tool a reason to stop", func(t *testing.T) {
		// when
		result := Build(&BuilderRequest{})
		// then
		assert.Contains(t, result, "never a precondition")
	})

	t.Run("carries a tool's instruction under its name", func(t *testing.T) {
		// given
		request := BuilderRequest{Tools: []mcp.Instruction{{Name: "file", Text: "reach only the folder /srv/analysis"}}}
		// when
		result := Build(&request)
		// then
		assert.Contains(t, result, `<tool-instructions name="file">`)
		assert.Contains(t, result, "reach only the folder /srv/analysis")
	})

	t.Run("orders the tool instructions by name", func(t *testing.T) {
		// given
		request := BuilderRequest{Tools: []mcp.Instruction{
			{Name: "web", Text: "third"},
			{Name: "date", Text: "first"},
			{Name: "file", Text: "second"},
		}}
		// when
		result := Build(&request)
		// then
		assert.Less(t, strings.Index(result, `name="date"`), strings.Index(result, `name="file"`))
		assert.Less(t, strings.Index(result, `name="file"`), strings.Index(result, `name="web"`))
	})

	t.Run("holds no instruction when no pack and no server has one", func(t *testing.T) {
		// when
		result := Build(&BuilderRequest{})
		// then
		assert.Contains(t, result, "<tools>")
		assert.NotContains(t, result, "<tool-instructions")
	})
}
