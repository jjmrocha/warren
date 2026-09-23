package setup

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jjmrocha/ai-toolkit/llm"
	"github.com/jjmrocha/warren/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testModel       = "z-ai/glm-5.3-flash"
	testOpenRouter  = "openrouter"
	testAnthropic   = "anthropic"
	testOllama      = "ollama"
	testOllamaModel = "qwen3"
)

func TestAskConfig(t *testing.T) {
	t.Run("collects every answer", func(t *testing.T) {
		// given
		in := bufio.NewReader(strings.NewReader("openrouter\n" + testModel + "\n"))
		// when
		result, err := askConfig(in, &strings.Builder{})
		// then
		require.NoError(t, err)
		assert.Equal(t, answers{provider: testOpenRouter, model: testModel}, result)
	})

	t.Run("opens with what setup is about to do", func(t *testing.T) {
		// given
		workDir(t)
		t.Setenv("HOME", t.TempDir())

		var out strings.Builder

		in := bufio.NewReader(strings.NewReader(testOllama + "\n" + testOllamaModel + "\n"))
		// when
		_, err := askConfig(in, &out)
		// then
		require.NoError(t, err)

		dir, err := os.Getwd()
		require.NoError(t, err)
		assert.Contains(t, out.String(), filepath.Join(dir, config.FileName))
		assert.Contains(t, out.String(), filepath.Join(os.Getenv("HOME"), ".claude", "skills"))
		assert.Contains(t, out.String(), "buffett-valuation")
		assert.Contains(t, out.String(), "company-research")
	})

	t.Run("lists the options it accepts", func(t *testing.T) {
		// given
		var out strings.Builder

		in := bufio.NewReader(strings.NewReader("ollama\nqwen3\n"))
		// when
		_, err := askConfig(in, &out)
		// then
		require.NoError(t, err)
		assert.Contains(t, out.String(), "anthropic, ollama, openrouter")
	})

	t.Run("asks again after an answer it cannot use", func(t *testing.T) {
		// given
		var out strings.Builder

		in := bufio.NewReader(strings.NewReader("openai\nanthropic\nbad name\nclaude-opus-5\n"))
		// when
		result, err := askConfig(in, &out)
		// then
		require.NoError(t, err)
		assert.Equal(t, answers{provider: testAnthropic, model: "claude-opus-5"}, result)
		assert.Contains(t, out.String(), `"openai" is not one of`)
		assert.Contains(t, out.String(), `"bad name" is not a valid answer`)
	})

	t.Run("reports answers it cannot read", func(t *testing.T) {
		// given
		in := bufio.NewReader(strings.NewReader("openrouter\n"))
		// when
		_, err := askConfig(in, &strings.Builder{})
		// then
		assert.ErrorIs(t, err, ErrNoAnswer)
	})
}

func TestKeyEnvFor(t *testing.T) {
	t.Run("names the variable the provider reads", func(t *testing.T) {
		testCases := []struct {
			name     string
			provider string
			expected string
		}{
			{name: testOpenRouter, provider: testOpenRouter, expected: "OPEN_ROUTER_KEY"},
			{name: testAnthropic, provider: testAnthropic, expected: "ANTHROPIC_API_KEY"},
			{name: testOllama, provider: testOllama, expected: ""},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				// when
				result := keyEnvFor(testCase.provider)
				// then
				assert.Equal(t, testCase.expected, result)
			})
		}
	})
}

func TestRenderConfig(t *testing.T) {
	t.Run("writes the config the answers describe", func(t *testing.T) {
		// given
		given := answers{provider: testOpenRouter, model: testModel}
		// when
		content, err := renderConfig(given)
		// then
		require.NoError(t, err)

		var result config.Config

		require.NoError(t, json.Unmarshal(content, &result))
		assert.Equal(t, "openrouter", result.LLM.Provider)
		assert.Equal(t, "OPEN_ROUTER_KEY", result.LLM.APIKeyEnv)
		assert.Equal(t, testModel, result.LLM.Model)
		assert.Equal(t, []string{testModel}, result.LLM.Models)
		assert.Equal(t, string(llm.EffortMedium), result.LLM.Effort)
		assert.Empty(t, result.Skills)
		assert.Empty(t, result.MCPsOn)
		assert.Equal(t, []string{"yfinance-mcp"}, names(result.MCPClients()))
		assert.Equal(t, 60*time.Second, result.MCPClients()[0].ToolCallTimeout)
	})

	t.Run("leaves out what the provider does not need", func(t *testing.T) {
		// given
		given := answers{provider: testOllama, model: testOllamaModel}
		// when
		content, err := renderConfig(given)
		// then
		require.NoError(t, err)
		assert.NotContains(t, string(content), "api-key-env")
		assert.NotContains(t, string(content), "base-url")
	})
}
