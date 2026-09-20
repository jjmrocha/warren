package engine

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/jjmrocha/warren/internal/config"
	"github.com/jjmrocha/warren/internal/prompt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testKeyEnv = "WARREN_TEST_KEY"

func testFile(skills string) string {
	return `{
  "llm": {
    "provider": "openrouter",
    "api-key-env": "` + testKeyEnv + `",
    "model": "z-ai/glm-5.3-flash",
    "effort": "medium"
  },
  "skills": ` + skills + `,
  "mcps": {},
  "mcps-on": []
}`
}

func testConfig(t *testing.T, content string) *config.Config {
	t.Helper()
	t.Setenv(testKeyEnv, "sk-test")

	dir := t.TempDir()
	t.Chdir(dir)

	path := filepath.Join(dir, config.FileName)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile(%s): %v", path, err)
	}

	cfg, err := config.Load()
	require.NoError(t, err)

	return cfg
}

func claudeSkills(t *testing.T, names ...string) {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)

	for _, name := range names {
		folder := filepath.Join(home, ".claude", "skills", name)
		if err := os.MkdirAll(folder, 0o750); err != nil {
			t.Fatalf("MkdirAll(%s): %v", folder, err)
		}

		content := "---\nname: " + name + "\ndescription: " + name + " skill\n---\n\nBody.\n"
		if err := os.WriteFile(filepath.Join(folder, "SKILL.md"), []byte(content), 0o600); err != nil {
			t.Fatalf("WriteFile(%s): %v", folder, err)
		}
	}
}

func TestNewSkillCollection(t *testing.T) {
	t.Run("reports every missing skill in one pass", func(t *testing.T) {
		// given
		cfg := testConfig(t, testFile(`[]`))
		claudeSkills(t)
		// when
		_, err := newSkillCollection(cfg)
		// then
		require.Error(t, err)

		for _, expected := range coreSkills {
			assert.Contains(t, err.Error(), expected)
		}
	})

	t.Run("names a missing extra skill beside the missing core ones", func(t *testing.T) {
		// given
		cfg := testConfig(t, testFile(`["removing-ai-tells"]`))
		claudeSkills(t)
		// when
		_, err := newSkillCollection(cfg)
		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "removing-ai-tells")
	})

	t.Run("loads every skill the config names", func(t *testing.T) {
		// given
		cfg := testConfig(t, testFile(`["removing-ai-tells"]`))
		wanted := slices.Concat(coreSkills, []string{"removing-ai-tells"})
		claudeSkills(t, wanted...)
		// when
		result, err := newSkillCollection(cfg)
		// then
		require.NoError(t, err)

		catalog := result.Catalog()
		for _, name := range wanted {
			assert.Contains(t, catalog, "<name>"+name+"</name>")
		}
	})

	t.Run("carries the two skills warren is built on", func(t *testing.T) {
		// given
		expected := []string{"buffett-valuation", "company-research"}
		// when
		result := coreSkills
		// then
		assert.Equal(t, expected, result)
	})

	t.Run("routes every core skill in the prompt", func(t *testing.T) {
		// given
		result := prompt.Build()
		// then
		for _, name := range coreSkills {
			assert.Contains(t, result, "| "+name)
		}
	})
}
