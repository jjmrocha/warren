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

const extraSkill = "removing-ai-tells"

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
		claudeSkills(t)
		// when
		_, err := newSkillCollection(&config.Config{})
		// then
		require.Error(t, err)

		for _, expected := range coreSkillNames() {
			assert.Contains(t, err.Error(), expected)
		}
	})

	t.Run("names a missing extra skill beside the missing core ones", func(t *testing.T) {
		// given
		claudeSkills(t)
		// when
		_, err := newSkillCollection(&config.Config{Skills: []string{extraSkill}})
		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), extraSkill)
	})

	t.Run("loads every skill the config names", func(t *testing.T) {
		// given
		wanted := slices.Concat(coreSkillNames(), []string{extraSkill})
		claudeSkills(t, wanted...)
		// when
		result, err := newSkillCollection(&config.Config{Skills: []string{extraSkill}})
		// then
		require.NoError(t, err)

		catalog := result.Catalog()
		for _, name := range wanted {
			assert.Contains(t, catalog, "<name>"+name+"</name>")
		}
	})

	t.Run("routes every core skill in the prompt", func(t *testing.T) {
		// given
		result := prompt.Build(nil)
		// then
		for _, name := range coreSkillNames() {
			assert.Contains(t, result, "| "+name)
		}
	})
}
