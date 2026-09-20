package setup

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jjmrocha/ai-toolkit/mcp"
	"github.com/jjmrocha/go-algo/fn"
	"github.com/jjmrocha/warren/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func workDir(t testing.TB) string {
	t.Helper()

	dir := t.TempDir()
	t.Chdir(dir)

	return dir
}

func answer(t *testing.T, content string) {
	t.Helper()

	path := filepath.Join(t.TempDir(), "answers")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile(%s): %v", path, err)
	}

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("Open(%s): %v", path, err)
	}

	stdin := os.Stdin
	os.Stdin = file

	t.Cleanup(func() {
		os.Stdin = stdin

		_ = file.Close()
	})
}

func names(configs []mcp.ClientConfig) []string {
	return fn.Map(configs, func(config mcp.ClientConfig) string {
		return config.Name
	})
}

func TestBuildIfNeed(t *testing.T) {
	t.Run("writes a config the working directory does not have", func(t *testing.T) {
		// given
		dir := workDir(t)
		answer(t, "ollama\nqwen3\n")
		// when
		err := BuildIfNeed()
		// then
		require.NoError(t, err)
		assert.FileExists(t, filepath.Join(dir, config.FileName))
	})

	t.Run("writes a config that loads back", func(t *testing.T) {
		testCases := []struct {
			name     string
			answers  string
			keyEnv   string
			provider string
		}{
			{name: testOpenRouter, answers: testOpenRouter + "\nz-ai/glm-5.3-flash\n", keyEnv: "OPEN_ROUTER_KEY", provider: testOpenRouter},
			{name: testAnthropic, answers: testAnthropic + "\nclaude-opus-5\n", keyEnv: "ANTHROPIC_API_KEY", provider: testAnthropic},
			{name: testOllama, answers: testOllama + "\nqwen3\n", provider: testOllama},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				// given
				if testCase.keyEnv != "" {
					t.Setenv(testCase.keyEnv, "sk-test")
				}

				workDir(t)
				answer(t, testCase.answers)
				require.NoError(t, BuildIfNeed())
				// when
				result, err := config.Load()
				// then
				require.NoError(t, err)
				assert.Equal(t, testCase.provider, result.LLM.Provider)
				assert.Equal(t, []string{"yfinance-mcp"}, names(result.MCPClients()))
				assert.Empty(t, result.MCPsOn)
			})
		}
	})

	t.Run("leaves an existing config alone", func(t *testing.T) {
		// given
		dir := workDir(t)
		path := filepath.Join(dir, config.FileName)
		require.NoError(t, os.WriteFile(path, []byte("{}"), 0o600))
		answer(t, "ollama\nqwen3\n")
		// when
		err := BuildIfNeed()
		// then
		require.NoError(t, err)

		content, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, "{}", string(content))
	})

	t.Run("stops on an answer it cannot read", func(t *testing.T) {
		// given
		dir := workDir(t)
		answer(t, "")
		// when
		err := BuildIfNeed()
		// then
		assert.ErrorIs(t, err, ErrNoAnswer)
		assert.NoFileExists(t, filepath.Join(dir, config.FileName))
	})
}
