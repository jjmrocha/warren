package engine

import (
	"testing"

	"github.com/jjmrocha/ai-toolkit/llm"
	"github.com/stretchr/testify/assert"
)

func fileCall(tool, path string) llm.ToolCall {
	return llm.ToolCall{ID: "call_1", Name: tool, Arguments: map[string]any{"path": path}}
}

func TestGuardFiles(t *testing.T) {
	t.Run("refuses a change to warren's own files", func(t *testing.T) {
		testCases := []struct {
			name string
			call llm.ToolCall
		}{
			{name: "write the config", call: fileCall("file_write", "warren.json")},
			{name: "edit the config", call: fileCall("file_edit", "warren.json")},
			{name: "delete the config", call: fileCall("file_delete", "warren.json")},
			{name: "config behind a dot", call: fileCall("file_write", "./warren.json")},
			{name: "config behind a detour", call: fileCall("file_write", "notes/../warren.json")},
			{name: "config in another case", call: fileCall("file_write", "Warren.JSON")},
			{name: "session file", call: fileCall("file_write", "sessions/4th2w29y76ozr1zthqah9dbnm.json")},
			{name: "sessions folder", call: fileCall("file_delete", "sessions")},
			{name: "sessions in another case", call: fileCall("file_edit", "Sessions/x.json")},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				// when
				result := guardFiles(t.Context(), testCase.call)
				// then
				assert.ErrorIs(t, result, ErrProtectedPath)
			})
		}
	})

	t.Run("lets every other call through", func(t *testing.T) {
		testCases := []struct {
			name string
			call llm.ToolCall
		}{
			{name: "read the config", call: fileCall("file_read", "warren.json")},
			{name: "read a session", call: fileCall("file_read", "sessions/x.json")},
			{name: "write a report", call: fileCall("file_write", "acme-valuation.md")},
			{name: "write a nested config name", call: fileCall("file_write", "notes/warren.json")},
			{name: "write a look-alike folder", call: fileCall("file_write", "sessions-notes/x.md")},
			{name: "a tool that is not a file tool", call: llm.ToolCall{Name: "current_date"}},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				// when
				result := guardFiles(t.Context(), testCase.call)
				// then
				assert.NoError(t, result)
			})
		}
	})
}
