package setup

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/jjmrocha/ai-toolkit/llm"
	"github.com/jjmrocha/warren/internal/config"
)

var modelPattern = regexp.MustCompile(`^[a-zA-Z0-9._:/-]+$`)

var keyEnvs = map[string]string{
	string(llm.ProviderOpenRouter): "OPEN_ROUTER_KEY",
	string(llm.ProviderAnthropic):  "ANTHROPIC_API_KEY",
}

type answers struct {
	provider string
	model    string
}

func askConfig(in *bufio.Reader, out io.Writer) (answers, error) {
	if err := intro(out); err != nil {
		return answers{}, err
	}

	providers := slices.Sorted(config.Providers.Values())
	choices := strings.Join(providers, ", ")

	provider, err := ask(in, out, "Provider ["+choices+"]: ", func(answer string) bool {
		return slices.Contains(providers, answer)
	}, "is not one of: "+choices)
	if err != nil {
		return answers{}, err
	}

	model, err := ask(in, out, "Model: ", modelPattern.MatchString, "is not a valid answer")
	if err != nil {
		return answers{}, err
	}

	return answers{provider: provider, model: model}, nil
}

func renderConfig(given answers) ([]byte, error) {
	cfg := config.Config{
		LLM: config.LLM{
			Provider:  given.provider,
			APIKeyEnv: keyEnvs[given.provider],
			Model:     given.model,
			Models:    []string{given.model},
			Effort:    string(llm.EffortMedium),
		},
		Skills: []string{},
		MCPs: map[string]config.MCP{
			"yfinance-mcp": {
				Command: "uvx",
				Args:    []string{"yfmcp@latest"},
				Timeout: 60,
			},
		},
		MCPsOn: []string{},
	}

	return json.MarshalIndent(cfg, "", "  ")
}

func ask(in *bufio.Reader, out io.Writer, prompt string, valid func(string) bool, complaint string) (string, error) {
	for {
		answer, err := read(in, out, prompt)
		if err != nil {
			return "", err
		}

		if valid(answer) {
			return answer, nil
		}

		_, _ = fmt.Fprintf(out, "%q %s\n", answer, complaint)
	}
}

func read(in *bufio.Reader, out io.Writer, prompt string) (string, error) {
	_, _ = fmt.Fprint(out, prompt)

	line, err := in.ReadString('\n')
	if errors.Is(err, io.EOF) && line == "" {
		return "", ErrNoAnswer
	}

	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	return strings.TrimSpace(line), nil
}

func intro(out io.Writer) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(out, `warren — The AI Financial Analyst.

First run in this folder: two questions, saved to

  %s

which you can edit later. warren also loads buffett-valuation and company-research from

  %s

and will not start without them — github.com/jjmrocha/investing-skills

`, filepath.Join(dir, config.FileName), filepath.Join(home, ".claude", "skills"))

	return err
}
