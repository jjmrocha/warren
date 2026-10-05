package engine

import (
	"testing"

	"github.com/jjmrocha/ai-toolkit/tools"
	"github.com/jjmrocha/warren/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func statusNames(t *testing.T, cfg *config.Config) []string {
	t.Helper()

	mng := newMCPManager(tools.NewToolBox(), cfg)
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
		cfg := &config.Config{MCPs: map[string]config.MCP{
			"yfinance-mcp": {Command: "uvx", Args: []string{"yfmcp@latest"}, Timeout: 60},
			"context7":     {Command: "npx", Args: []string{"-y", "@upstash/context7-mcp"}},
		}}
		// when
		result := statusNames(t, cfg)
		// then
		assert.ElementsMatch(t, []string{"yfinance-mcp", "context7"}, result)
	})

	t.Run("registers nothing when the config has no servers", func(t *testing.T) {
		// when
		result := statusNames(t, &config.Config{})
		// then
		assert.Empty(t, result)
	})
}

func TestStartMCPs(t *testing.T) {
	t.Run("carries on when a boot server fails to start", func(t *testing.T) {
		// given
		cfg := &config.Config{
			MCPs:   map[string]config.MCP{"broken": {Command: "definitely-not-a-binary"}},
			MCPsOn: []string{"broken"},
		}

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
