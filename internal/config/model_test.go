package config

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/jjmrocha/ai-toolkit/llm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func configWith(overrides ...string) string {
	fields := map[string]string{
		"llm":     `{"provider": "openrouter", "api-key-env": "` + testKeyEnv + `", "model": "` + testModel + `", "effort": "medium"}`,
		"skills":  `[]`,
		"mcps":    `{}`,
		"mcps-on": `[]`,
	}

	for _, override := range overrides {
		key, value, _ := strings.Cut(override, ":")
		fields[strings.Trim(strings.TrimSpace(key), `"`)] = strings.TrimSpace(value)
	}

	parts := make([]string, 0, len(fields))
	for key, value := range fields {
		parts = append(parts, `"`+key+`": `+value)
	}

	return "{" + strings.Join(parts, ",") + "}"
}

func mustLoad(t *testing.T, content string) *Config {
	t.Helper()
	t.Setenv(testKeyEnv, "sk-test")

	dir := workDir(t)
	writeConfig(t, dir, content)

	config, err := Load()
	require.NoError(t, err)

	return config
}

func TestProviders(t *testing.T) {
	t.Run("holds every provider warren accepts", func(t *testing.T) {
		// given
		expected := []string{"anthropic", "ollama", "openrouter"}
		// when
		result := slices.Sorted(Providers.Values())
		// then
		assert.Equal(t, expected, result)
	})
}

func TestLLMConfig(t *testing.T) {
	t.Run("carries every value the model needs", func(t *testing.T) {
		// given
		config := mustLoad(t, validConfig())
		// when
		result := config.LLMConfig()
		// then
		assert.Equal(t, llm.ProviderOpenRouter, result.Provider)
		assert.Equal(t, testModel, result.Model)
		assert.Equal(t, []string{testModel, "deepseek/deepseek-v4-pro"}, result.Models)
		assert.Equal(t, llm.EffortMedium, result.Effort)
		assert.Equal(t, "sk-test", result.APIKey)
		assert.Empty(t, result.BaseURL)
	})
}

func TestMCPClients(t *testing.T) {
	t.Run("converts every entry into a client config", func(t *testing.T) {
		// given
		content := configWith(`"mcps": {"yfinance-mcp": {"command": "uvx", "args": ["yfmcp@latest"], "env": ["HOME"], "timeout": 90}}`)
		config := mustLoad(t, content)
		// when
		result := config.MCPClients()
		// then
		require.Len(t, result, 1)
		assert.Equal(t, "yfinance-mcp", result[0].Name)
		assert.Equal(t, "uvx", result[0].Command)
		assert.Equal(t, []string{"yfmcp@latest"}, result[0].Args)
		assert.Equal(t, []string{"HOME"}, result[0].InheritEnv)
		assert.Equal(t, 90*time.Second, result[0].ToolCallTimeout)
	})

	t.Run("leaves the timeout zero when none is set", func(t *testing.T) {
		// given
		content := configWith(`"mcps": {"yfinance-mcp": {"command": "uvx", "args": ["yfmcp@latest"]}}`)
		config := mustLoad(t, content)
		// when
		result := config.MCPClients()
		// then
		require.Len(t, result, 1)
		assert.Zero(t, result[0].ToolCallTimeout)
	})

	t.Run("returns nothing when the section is empty", func(t *testing.T) {
		// given
		config := mustLoad(t, configWith())
		// when
		result := config.MCPClients()
		// then
		assert.Empty(t, result)
	})
}
