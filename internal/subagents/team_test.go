package subagents

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jjmrocha/ai-chat/chat"
	"github.com/jjmrocha/ai-toolkit/agent"
	"github.com/jjmrocha/ai-toolkit/llm"
	"github.com/jjmrocha/ai-toolkit/subagent"
	"github.com/jjmrocha/ai-toolkit/tools"
	"github.com/jjmrocha/warren/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testWebReaderModel     = "deepseek/deepseek-v4-flash"
	testFinancialDataModel = "z-ai/glm-5.3"
	answerReply            = `{"message":{"content":"the answer"},"done":true,"prompt_eval_count":20,"eval_count":10}`
	modelInfo              = `{"model_info":{"general.architecture":"test","test.context_length":8192}}`
)

func newConfig(webReaderOn, financialDataOn bool) *config.Config {
	return &config.Config{
		LLM: config.LLM{Provider: string(llm.ProviderOllama), Model: "main", Effort: string(llm.EffortMedium)},
		Subagents: config.Subagents{
			WebReader:     config.Subagent{Model: testWebReaderModel, Effort: string(llm.EffortLow), Enabled: webReaderOn},
			FinancialData: config.Subagent{Model: testFinancialDataModel, Effort: string(llm.EffortMax), Enabled: financialDataOn},
		},
	}
}

func newModel(t *testing.T, onChat func()) *llm.LLM {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path != "/api/chat" {
			_, _ = io.WriteString(w, modelInfo)
			return
		}

		onChat()

		_, _ = io.WriteString(w, answerReply)
	}))
	t.Cleanup(server.Close)

	model, err := llm.New(llm.Config{Provider: llm.ProviderOllama, Model: "test", BaseURL: server.URL})
	require.NoError(t, err)

	return model
}

func newMainAgent(t *testing.T) *agent.Agent {
	t.Helper()

	mainAgent, err := agent.New(agent.Config{}, newModel(t, func() {}))
	require.NoError(t, err)

	return mainAgent
}

type fakeSubagent struct {
	configs []llm.Config
	closed  int
	onChat  func()
}

func (f *fakeSubagent) register(t *testing.T, name, description string, err error) func(context.Context, *subagent.Manager, llm.Config) (func(), error) {
	t.Helper()

	return func(_ context.Context, manager *subagent.Manager, llmConfig llm.Config) (func(), error) {
		f.configs = append(f.configs, llmConfig)

		if err != nil {
			return nil, err
		}

		child, newErr := agent.New(agent.Config{}, newModel(t, func() { f.onChat() }))
		require.NoError(t, newErr)

		fakeConfig := subagent.Config{Name: name, Description: description}
		require.NoError(t, manager.Register(fakeConfig, child, nil, nil))

		return func() { f.closed++ }, nil
	}
}

func installFakeWebReader(t *testing.T, err error) *fakeSubagent {
	t.Helper()

	fake := &fakeSubagent{onChat: func() {}}

	previous := registerWebReader
	t.Cleanup(func() { registerWebReader = previous })
	registerWebReader = fake.register(t, "web-reader", "Reads the web", err)

	return fake
}

func installFakeFinancialData(t *testing.T, err error) *fakeSubagent {
	t.Helper()

	fake := &fakeSubagent{onChat: func() {}}

	previous := registerFinancialData
	t.Cleanup(func() { registerFinancialData = previous })
	registerFinancialData = fake.register(t, "financial-data", "Reads Yahoo", err)

	return fake
}

type fakeProgress struct {
	texts []string
}

func (f *fakeProgress) SetText(text string) { f.texts = append(f.texts, text) }

func TestStart(t *testing.T) {
	t.Run("adds no subagent tool when no subagent is enabled", func(t *testing.T) {
		// given
		installFakeFinancialData(t, nil)
		fake := installFakeWebReader(t, nil)
		toolBox := tools.NewToolBox()
		line := &fakeProgress{}
		// when
		team, err := Start(context.Background(), toolBox, newConfig(false, false), line)
		// then
		require.NoError(t, err)
		require.NotNil(t, team)
		_, found := toolBox.Tool("subagent")
		assert.False(t, found)
		assert.Empty(t, fake.configs)
		assert.Empty(t, line.texts)
	})

	t.Run("names the web-reader before it starts", func(t *testing.T) {
		// given
		installFakeWebReader(t, nil)
		line := &fakeProgress{}
		expected := []string{"Starting subagent web-reader…"}
		// when
		team, err := Start(context.Background(), tools.NewToolBox(), newConfig(true, false), line)
		// then
		require.NoError(t, err)
		t.Cleanup(team.Close)
		assert.Equal(t, expected, line.texts)
	})

	t.Run("adds the subagent tool listing the web-reader when it is enabled", func(t *testing.T) {
		// given
		installFakeWebReader(t, nil)
		toolBox := tools.NewToolBox()
		// when
		team, err := Start(context.Background(), toolBox, newConfig(true, false), &fakeProgress{})
		// then
		require.NoError(t, err)
		t.Cleanup(team.Close)

		result, found := toolBox.Tool("subagent")
		require.True(t, found)
		assert.Contains(t, result.Description, "- web-reader: Reads the web")
	})

	t.Run("gives the web-reader the main settings with its own model and effort", func(t *testing.T) {
		// given
		fake := installFakeWebReader(t, nil)
		cfg := newConfig(true, false)
		expected := cfg.SubagentLLMConfig(cfg.Subagents.WebReader)
		// when
		team, err := Start(context.Background(), tools.NewToolBox(), cfg, &fakeProgress{})
		// then
		require.NoError(t, err)
		t.Cleanup(team.Close)
		assert.Equal(t, []llm.Config{expected}, fake.configs)
	})

	t.Run("names the web-reader when it cannot be built", func(t *testing.T) {
		// given
		failure := errors.New("no npx")
		installFakeWebReader(t, failure)
		toolBox := tools.NewToolBox()
		// when
		_, err := Start(context.Background(), toolBox, newConfig(true, false), &fakeProgress{})
		// then
		require.ErrorIs(t, err, failure)
		assert.Contains(t, err.Error(), "web-reader: ")
		_, found := toolBox.Tool("subagent")
		assert.False(t, found)
	})
	t.Run("adds the subagent tool listing only the financial-data agent when it alone is enabled", func(t *testing.T) {
		// given
		webReader := installFakeWebReader(t, nil)
		installFakeFinancialData(t, nil)
		toolBox := tools.NewToolBox()
		// when
		team, err := Start(context.Background(), toolBox, newConfig(false, true), &fakeProgress{})
		// then
		require.NoError(t, err)
		t.Cleanup(team.Close)

		result, found := toolBox.Tool("subagent")
		require.True(t, found)
		assert.Contains(t, result.Description, "- financial-data: Reads Yahoo")
		assert.NotContains(t, result.Description, "- web-reader:")
		assert.Empty(t, webReader.configs)
	})

	t.Run("lists both subagents when both are enabled", func(t *testing.T) {
		// given
		installFakeWebReader(t, nil)
		installFakeFinancialData(t, nil)
		toolBox := tools.NewToolBox()
		// when
		team, err := Start(context.Background(), toolBox, newConfig(true, true), &fakeProgress{})
		// then
		require.NoError(t, err)
		t.Cleanup(team.Close)

		result, found := toolBox.Tool("subagent")
		require.True(t, found)
		assert.Contains(t, result.Description, "- web-reader: Reads the web")
		assert.Contains(t, result.Description, "- financial-data: Reads Yahoo")
	})

	t.Run("names each subagent before it starts", func(t *testing.T) {
		// given
		installFakeWebReader(t, nil)
		installFakeFinancialData(t, nil)
		line := &fakeProgress{}
		expected := []string{"Starting subagent web-reader…", "Starting subagent financial-data…"}
		// when
		team, err := Start(context.Background(), tools.NewToolBox(), newConfig(true, true), line)
		// then
		require.NoError(t, err)
		t.Cleanup(team.Close)
		assert.Equal(t, expected, line.texts)
	})

	t.Run("gives the financial-data agent the main settings with its own model and effort", func(t *testing.T) {
		// given
		installFakeWebReader(t, nil)
		fake := installFakeFinancialData(t, nil)
		cfg := newConfig(true, true)
		expected := cfg.SubagentLLMConfig(cfg.Subagents.FinancialData)
		// when
		team, err := Start(context.Background(), tools.NewToolBox(), cfg, &fakeProgress{})
		// then
		require.NoError(t, err)
		t.Cleanup(team.Close)
		assert.Equal(t, []llm.Config{expected}, fake.configs)
	})

	t.Run("closes the web-reader and names the financial-data agent when it cannot be built", func(t *testing.T) {
		// given
		webReader := installFakeWebReader(t, nil)
		failure := errors.New("no uvx")
		installFakeFinancialData(t, failure)
		toolBox := tools.NewToolBox()
		// when
		_, err := Start(context.Background(), toolBox, newConfig(true, true), &fakeProgress{})
		// then
		require.ErrorIs(t, err, failure)
		assert.Contains(t, err.Error(), "financial-data: ")
		assert.Equal(t, 1, webReader.closed)
		_, found := toolBox.Tool("subagent")
		assert.False(t, found)
	})
}

func TestTeamChatOption(t *testing.T) {
	t.Run("shows the web-reader's progress in the chat while it runs", func(t *testing.T) {
		// given
		fake := installFakeWebReader(t, nil)
		toolBox := tools.NewToolBox()

		team, err := Start(context.Background(), toolBox, newConfig(true, false), &fakeProgress{})
		require.NoError(t, err)
		t.Cleanup(team.Close)

		chatCore := chat.New("WARREN", newMainAgent(t), team.ChatOption())

		var seen []string

		fake.onChat = func() { seen = append(seen, chatCore.Progress().Subagent) }
		call := llm.ToolCall{ID: "1", Name: "subagent", Arguments: map[string]any{"agent": "web-reader", "prompt": "Ferrari FY2025 revenue"}}
		// when
		_, err = toolBox.Execute(context.Background(), call)
		// then
		require.NoError(t, err)
		require.Len(t, seen, 1)
		assert.NotEmpty(t, seen[0])
		assert.Empty(t, chatCore.Progress().Subagent)
	})

	t.Run("leaves the chat without a subagent line when the web-reader is disabled", func(t *testing.T) {
		// given
		installFakeWebReader(t, nil)

		team, err := Start(context.Background(), tools.NewToolBox(), newConfig(false, false), &fakeProgress{})
		require.NoError(t, err)
		// when
		chatCore := chat.New("WARREN", newMainAgent(t), team.ChatOption())
		// then
		assert.Empty(t, chatCore.Progress().Subagent)
	})
}

func TestTeamClose(t *testing.T) {
	t.Run("removes the subagent tool and closes the web-reader", func(t *testing.T) {
		// given
		fake := installFakeWebReader(t, nil)
		toolBox := tools.NewToolBox()

		team, err := Start(context.Background(), toolBox, newConfig(true, false), &fakeProgress{})
		require.NoError(t, err)
		// when
		team.Close()
		// then
		_, found := toolBox.Tool("subagent")
		assert.False(t, found)
		assert.Equal(t, 1, fake.closed)
	})

	t.Run("does nothing when called again", func(t *testing.T) {
		// given
		fake := installFakeWebReader(t, nil)

		team, err := Start(context.Background(), tools.NewToolBox(), newConfig(true, false), &fakeProgress{})
		require.NoError(t, err)

		team.Close()
		// when
		team.Close()
		// then
		assert.Equal(t, 1, fake.closed)
	})
}
