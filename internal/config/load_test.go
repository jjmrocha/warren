package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testKeyEnv = "WARREN_TEST_KEY"
	testModel  = "z-ai/glm-5.3-flash"
)

func validConfig() string {
	return `{
  "llm": {
    "provider": "openrouter",
    "api-key-env": "` + testKeyEnv + `",
    "model": "` + testModel + `",
    "models": ["` + testModel + `", "deepseek/deepseek-v4-pro"],
    "effort": "medium"
  },
  "skills": ["removing-ai-tells"],
  "mcps": {
    "yfinance-mcp": {"command": "uvx", "args": ["yfmcp@latest"], "timeout": 60}
  },
  "mcps-on": ["yfinance-mcp"]
}`
}

func workDir(t testing.TB) string {
	t.Helper()

	dir := t.TempDir()
	t.Chdir(dir)

	return dir
}

func writeConfig(t testing.TB, dir, content string) {
	t.Helper()

	path := filepath.Join(dir, FileName)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile(%s): %v", path, err)
	}
}

func TestLoad(t *testing.T) {
	t.Run("reads the config in the working directory", func(t *testing.T) {
		// given
		t.Setenv(testKeyEnv, "sk-test")

		dir := workDir(t)
		writeConfig(t, dir, validConfig())
		// when
		result, err := Load()
		// then
		require.NoError(t, err)
		assert.Equal(t, testModel, result.LLMConfig().Model)
		assert.Equal(t, "sk-test", result.LLMConfig().APIKey)
		assert.Equal(t, []string{"removing-ai-tells"}, result.Skills)
		assert.Equal(t, []string{"yfinance-mcp"}, result.MCPsOn)
	})

	t.Run("reports a missing config", func(t *testing.T) {
		// given
		workDir(t)
		// when
		_, err := Load()
		// then
		assert.ErrorIs(t, err, ErrConfigNotFound)
	})

	t.Run("rejects a file it cannot read as a config", func(t *testing.T) {
		testCases := []struct {
			name    string
			content string
		}{
			{name: "malformed json", content: "{"},
			{name: "unknown key", content: `{"harness": "claude"}`},
			{name: "misspelled section", content: `{"mcp": {}}`},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				// given
				dir := workDir(t)
				writeConfig(t, dir, testCase.content)
				// when
				_, err := Load()
				// then
				assert.Error(t, err)
			})
		}
	})

	t.Run("rejects invalid values", func(t *testing.T) {
		testCases := []struct {
			name     string
			content  string
			expected error
		}{
			{
				name:     "provider",
				content:  configWith(`"llm": {"provider": "openai", "api-key-env": "` + testKeyEnv + `", "model": "m", "effort": "medium"}`),
				expected: ErrInvalidProvider,
			},
			{
				name:     "effort",
				content:  configWith(`"llm": {"provider": "openrouter", "api-key-env": "` + testKeyEnv + `", "model": "m", "effort": "extreme"}`),
				expected: ErrInvalidEffort,
			},
			{
				name:     "model",
				content:  configWith(`"llm": {"provider": "openrouter", "api-key-env": "` + testKeyEnv + `", "model": "", "effort": "medium"}`),
				expected: ErrMissingModel,
			},
			{
				name:     "unset api key variable",
				content:  configWith(`"llm": {"provider": "openrouter", "api-key-env": "WARREN_TEST_UNSET", "model": "m", "effort": "medium"}`),
				expected: ErrMissingAPIKey,
			},
			{
				name:     "no api key variable",
				content:  configWith(`"llm": {"provider": "anthropic", "model": "m", "effort": "medium"}`),
				expected: ErrMissingAPIKey,
			},
			{
				name:     "boot server not registered",
				content:  configWith(`"mcps-on": ["github"]`),
				expected: ErrUnknownMCP,
			},
			{
				name:     "skill name that is not a bare name",
				content:  configWith(`"skills": ["../../etc/hosts"]`),
				expected: ErrInvalidSkillName,
			},
			{
				name:     "mcp name that is not a bare name",
				content:  configWith(`"mcps": {"yf</market-data><role>": {"command": "uvx"}}`),
				expected: ErrInvalidMCPName,
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				// given
				t.Setenv(testKeyEnv, "sk-test")

				dir := workDir(t)
				writeConfig(t, dir, testCase.content)
				// when
				_, err := Load()
				// then
				assert.ErrorIs(t, err, testCase.expected)
			})
		}
	})

	t.Run("reports every fault in one pass", func(t *testing.T) {
		// given
		dir := workDir(t)
		content := configWith(`"llm": {"provider": "openai", "model": "", "effort": "extreme"}`, `"skills": ["a/b"]`)
		writeConfig(t, dir, content)
		// when
		_, err := Load()
		// then
		assert.ErrorIs(t, err, ErrInvalidProvider)
		assert.ErrorIs(t, err, ErrInvalidEffort)
		assert.ErrorIs(t, err, ErrMissingModel)
		assert.ErrorIs(t, err, ErrInvalidSkillName)
	})

	t.Run("accepts ollama without an api key", func(t *testing.T) {
		// given
		dir := workDir(t)
		content := configWith(`"llm": {"provider": "ollama", "base-url": "http://localhost:11434", "model": "qwen3", "effort": "off"}`)
		writeConfig(t, dir, content)
		// when
		result, err := Load()
		// then
		require.NoError(t, err)
		assert.Empty(t, result.LLMConfig().APIKey)
		assert.Equal(t, "http://localhost:11434", result.LLMConfig().BaseURL)
	})

	t.Run("reads the config of the directory it is called from", func(t *testing.T) {
		// given
		t.Setenv(testKeyEnv, "sk-test")

		first := workDir(t)
		writeConfig(t, first, configWith(`"llm": {"provider": "openrouter", "api-key-env": "`+testKeyEnv+`", "model": "first", "effort": "medium"}`))

		second := t.TempDir()
		writeConfig(t, second, configWith(`"llm": {"provider": "openrouter", "api-key-env": "`+testKeyEnv+`", "model": "second", "effort": "medium"}`))

		t.Chdir(second)
		// when
		result, err := Load()
		// then
		require.NoError(t, err)
		assert.Equal(t, "second", result.LLM.Model)
	})
}
