package engine

import (
	"testing"

	"github.com/jjmrocha/ai-toolkit/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func statusNames(t *testing.T, content string) []string {
	t.Helper()

	mng := newMCPManager(tools.NewToolBox(), testConfig(t, content))
	t.Cleanup(mng.Close)

	names := make([]string, 0)
	for _, status := range mng.Status() {
		names = append(names, status.Name)
		assert.False(t, status.Active)
	}

	return names
}

func TestNewMCPManager(t *testing.T) {
	t.Run("registers every server the config names", func(t *testing.T) {
		// given
		content := `{
  "llm": {"provider": "openrouter", "api-key-env": "` + testKeyEnv + `", "model": "m", "effort": "medium"},
  "skills": [],
  "mcps": {
    "yfinance-mcp": {"command": "uvx", "args": ["yfmcp@latest"], "timeout": 60},
    "context7": {"command": "npx", "args": ["-y", "@upstash/context7-mcp"]}
  },
  "mcps-on": []
}`
		// when
		result := statusNames(t, content)
		// then
		assert.ElementsMatch(t, []string{"yfinance-mcp", "context7"}, result)
	})

	t.Run("registers nothing when the config has no servers", func(t *testing.T) {
		// given
		content := testFile(`[]`)
		// when
		result := statusNames(t, content)
		// then
		assert.Empty(t, result)
	})
}

func TestStartMCPs(t *testing.T) {
	t.Run("carries on when a boot server fails to start", func(t *testing.T) {
		// given
		content := `{
  "llm": {"provider": "openrouter", "api-key-env": "` + testKeyEnv + `", "model": "m", "effort": "medium"},
  "skills": [],
  "mcps": {"broken": {"command": "definitely-not-a-binary"}},
  "mcps-on": ["broken"]
}`
		cfg := testConfig(t, content)

		mng := newMCPManager(tools.NewToolBox(), cfg)
		t.Cleanup(mng.Close)
		// when
		startMCPs(t.Context(), mng, cfg)
		// then
		result := mng.Status()
		require.Len(t, result, 1)
		assert.Equal(t, "broken", result[0].Name)
		assert.False(t, result[0].Active)
	})
}
