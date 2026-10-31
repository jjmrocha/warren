package financialdata

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/jjmrocha/ai-toolkit/llm"
	"github.com/jjmrocha/ai-toolkit/mcp"
	"github.com/jjmrocha/ai-toolkit/packs"
	"github.com/jjmrocha/ai-toolkit/subagent"
	"github.com/jjmrocha/ai-toolkit/tools"
	"github.com/jjmrocha/go-algo/fn"
	"github.com/jjmrocha/warren/internal/subagents/financialdata/prompt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	answerReply     = `{"message":{"content":"the answer"},"done":true,"prompt_eval_count":20,"eval_count":10}`
	modelInfo       = `{"model_info":{"general.architecture":"test","test.context_length":8192}}`
	testFinanceTool = "yfinance__yfinance_get_ticker_info"
	financeManual   = "Quotes may be delayed."
)

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatTool struct {
	Function struct {
		Name string `json:"name"`
	} `json:"function"`
}

type chatRequest struct {
	Messages []chatMessage `json:"messages"`
	Tools    []chatTool    `json:"tools"`
}

func (r chatRequest) systemMessage() string {
	for _, message := range r.Messages {
		if message.Role == "system" {
			return message.Content
		}
	}

	return ""
}

func (r chatRequest) toolNames() []string {
	return fn.Map(r.Tools, func(tool chatTool) string { return tool.Function.Name })
}

type modelServer struct {
	mu       sync.Mutex
	requests []chatRequest
	url      string
}

func (s *modelServer) record(request chatRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.requests = append(s.requests, request)
}

func (s *modelServer) firstRequest(t *testing.T) chatRequest {
	t.Helper()

	s.mu.Lock()
	defer s.mu.Unlock()

	require.NotEmpty(t, s.requests)

	return s.requests[0]
}

func newModelServer(t *testing.T) *modelServer {
	t.Helper()

	recorder := &modelServer{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path != "/api/chat" {
			_, _ = io.WriteString(w, modelInfo)
			return
		}

		var request chatRequest

		_ = json.NewDecoder(r.Body).Decode(&request)
		recorder.record(request)

		_, _ = io.WriteString(w, answerReply)
	}))
	t.Cleanup(server.Close)

	recorder.url = server.URL

	return recorder
}

type fakePack struct {
	closed bool
}

func (p *fakePack) Close() error {
	p.closed = true
	return nil
}

func (p *fakePack) Instructions(_ context.Context) *mcp.Instruction {
	return &mcp.Instruction{Name: "yfinance", Text: financeManual}
}

func fakeFinanceTools(t *testing.T, err error) *fakePack {
	t.Helper()

	pack := &fakePack{}

	previous := financeTools
	t.Cleanup(func() { financeTools = previous })
	financeTools = func(_ context.Context, m *tools.ToolBox) (packs.ToolPack, error) {
		if err != nil {
			return nil, err
		}

		tool := llm.Tool{Name: testFinanceTool, Description: testFinanceTool}
		if addErr := m.Add(tool, func(context.Context, map[string]any) (string, error) { return "", nil }); addErr != nil {
			return nil, addErr
		}

		return pack, nil
	}

	return pack
}

func newLLMConfig(server *modelServer) llm.Config {
	return llm.Config{Provider: llm.ProviderOllama, Model: "test", BaseURL: server.url}
}

func TestRegister(t *testing.T) {
	t.Run("registers the financial-data agent with its description", func(t *testing.T) {
		// given
		server := newModelServer(t)
		fakeFinanceTools(t, nil)
		manager := subagent.NewManager()
		expected := []subagent.Info{{Name: "financial-data", Description: description}}
		// when
		closeFn, err := Register(context.Background(), manager, newLLMConfig(server))
		// then
		require.NoError(t, err)
		t.Cleanup(closeFn)
		assert.Equal(t, expected, manager.Agents())
	})

	t.Run("offers only the finance tools", func(t *testing.T) {
		// given
		server := newModelServer(t)
		fakeFinanceTools(t, nil)
		manager := subagent.NewManager()

		closeFn, err := Register(context.Background(), manager, newLLMConfig(server))
		require.NoError(t, err)
		t.Cleanup(closeFn)
		// when
		result, err := manager.Run(context.Background(), "financial-data", "ASML.AS income statement, FY2023 to FY2025")
		// then
		require.NoError(t, err)
		assert.Equal(t, "the answer", result)
		assert.Equal(t, []string{testFinanceTool}, server.firstRequest(t).toolNames())
	})

	t.Run("builds the system prompt from the finance tools' instructions", func(t *testing.T) {
		// given
		server := newModelServer(t)
		fakeFinanceTools(t, nil)
		manager := subagent.NewManager()

		closeFn, err := Register(context.Background(), manager, newLLMConfig(server))
		require.NoError(t, err)
		t.Cleanup(closeFn)

		expected := prompt.Build(prompt.Request{
			ToolInstructions: []mcp.Instruction{{Name: "yfinance", Text: financeManual}},
		})
		// when
		_, err = manager.Run(context.Background(), "financial-data", "ASML.AS income statement, FY2023 to FY2025")
		// then
		require.NoError(t, err)
		assert.Equal(t, expected, server.firstRequest(t).systemMessage())
	})

	t.Run("closes the finance tools through the function it returns", func(t *testing.T) {
		// given
		server := newModelServer(t)
		finance := fakeFinanceTools(t, nil)
		manager := subagent.NewManager()

		closeFn, err := Register(context.Background(), manager, newLLMConfig(server))
		require.NoError(t, err)
		// when
		closeFn()
		// then
		assert.True(t, finance.closed)
	})

	t.Run("registers nothing when the finance tools fail to start", func(t *testing.T) {
		// given
		server := newModelServer(t)
		failure := errors.New("no uvx")
		fakeFinanceTools(t, failure)
		manager := subagent.NewManager()
		// when
		_, err := Register(context.Background(), manager, newLLMConfig(server))
		// then
		require.ErrorIs(t, err, failure)
		assert.Contains(t, err.Error(), "finance tools: ")
		assert.Empty(t, manager.Agents())
	})
}
