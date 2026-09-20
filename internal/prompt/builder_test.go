package prompt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	t.Run("routes research ahead of valuation", func(t *testing.T) {
		// when
		result := Build()
		// then
		assert.Less(t, strings.Index(result, "company-research"), strings.Index(result, "buffett-valuation"))
	})

	t.Run("names the format and the shape of a report it is asked to write", func(t *testing.T) {
		// when
		result := Build()
		// then
		assert.Contains(t, result, "Markdown")
		assert.Contains(t, result, "subfolder")
	})

	t.Run("sends the model to the date tool instead of its own sense of the date", func(t *testing.T) {
		// when
		result := Build()
		// then
		assert.Contains(t, result, "current_date")
	})

	t.Run("tells the model to hand a skill's script an absolute path", func(t *testing.T) {
		// when
		result := Build()
		// then
		assert.Contains(t, result, "file_workdir")
	})

	t.Run("tells the model to remove a working file it no longer needs", func(t *testing.T) {
		// when
		result := Build()
		// then
		assert.Contains(t, result, "file_delete")
	})

	t.Run("names the tool that runs a skill's scripts", func(t *testing.T) {
		// when
		result := Build()
		// then
		assert.Contains(t, result, "skill_execute_file")
	})

	t.Run("never makes a missing tool a reason to stop", func(t *testing.T) {
		// when
		result := Build()
		// then
		assert.Contains(t, result, "never a precondition")
	})
}
