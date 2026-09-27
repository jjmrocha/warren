package engine

import (
	"context"
	"errors"
	"testing"

	"github.com/jjmrocha/ai-toolkit/mcp"
	"github.com/jjmrocha/ai-toolkit/packs"
	"github.com/jjmrocha/ai-toolkit/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type silentPack struct{}

func (silentPack) Close() error { return nil }

func (silentPack) Instructions(context.Context) (*mcp.Instruction, error) { return nil, nil }

type failingPack struct{}

func (failingPack) Close() error { return nil }

func (failingPack) Instructions(context.Context) (*mcp.Instruction, error) {
	return &mcp.Instruction{Name: "failing", Text: "never used"}, errors.New("pack gone")
}

func emptyManager(t *testing.T) *mcp.Manager {
	t.Helper()

	content := `{
  "llm": {"provider": "openrouter", "api-key-env": "` + testKeyEnv + `", "model": "m", "effort": "medium"},
  "skills": [],
  "mcps": {},
  "mcps-on": []
}`

	mng := newMCPManager(tools.NewToolBox(), testConfig(t, content))
	t.Cleanup(mng.Close)

	return mng
}

func instructionNames(instructions []mcp.Instruction) []string {
	names := make([]string, 0, len(instructions))
	for _, instruction := range instructions {
		names = append(names, instruction.Name)
	}

	return names
}

func TestToolInstructions(t *testing.T) {
	t.Run("collects the instruction of a pack that has one", func(t *testing.T) {
		// given
		ctx := context.Background()
		datePack, err := packs.DateTools(tools.NewToolBox())
		require.NoError(t, err)
		t.Cleanup(func() { _ = datePack.Close() })
		expected, err := datePack.Instructions(ctx)
		require.NoError(t, err)
		// when
		result := toolInstructions(ctx, []packs.ToolPack{datePack}, emptyManager(t))
		// then
		require.Len(t, result, 1)
		assert.Equal(t, *expected, result[0])
	})

	t.Run("skips a pack that has no instruction", func(t *testing.T) {
		// given
		ctx := context.Background()
		// when
		result := toolInstructions(ctx, []packs.ToolPack{silentPack{}}, emptyManager(t))
		// then
		assert.Empty(t, result)
	})

	t.Run("skips a pack whose instruction fails", func(t *testing.T) {
		// given
		ctx := context.Background()
		// when
		result := toolInstructions(ctx, []packs.ToolPack{failingPack{}}, emptyManager(t))
		// then
		assert.Empty(t, result)
	})

	t.Run("keeps the packs that work when one of them fails", func(t *testing.T) {
		// given
		ctx := context.Background()
		datePack, err := packs.DateTools(tools.NewToolBox())
		require.NoError(t, err)
		t.Cleanup(func() { _ = datePack.Close() })
		// when
		result := toolInstructions(ctx, []packs.ToolPack{failingPack{}, datePack, silentPack{}}, emptyManager(t))
		// then
		assert.Equal(t, []string{"date"}, instructionNames(result))
	})

	t.Run("returns nothing when there are no packs and no servers", func(t *testing.T) {
		// given
		ctx := context.Background()
		// when
		result := toolInstructions(ctx, nil, emptyManager(t))
		// then
		assert.Empty(t, result)
	})
}
