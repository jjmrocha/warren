package setup

import (
	"bufio"
	"encoding/json"
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
	var given answers

	if err := intro(out); err != nil {
		return given, err
	}

	provider, err := choose(in, out, "Provider", slices.Sorted(config.Providers.Values()))
	if err != nil {
		return given, err
	}

	given.provider = provider

	model, err := match(in, out, "Model", modelPattern)
	if err != nil {
		return given, err
	}

	given.model = model

	return given, nil
}

func keyEnvFor(provider string) string {
	return keyEnvs[provider]
}

func renderConfig(given answers) ([]byte, error) {
	cfg := config.Config{
		LLM: config.LLM{
			Provider:  given.provider,
			APIKeyEnv: keyEnvFor(given.provider),
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

func choose(in *bufio.Reader, out io.Writer, question string, options []string) (string, error) {
	prompt := fmt.Sprintf("%s [%s]: ", question, strings.Join(options, ", "))

	for {
		answer, err := read(in, out, prompt)
		if err != nil {
			return "", err
		}

		if slices.Contains(options, answer) {
			return answer, nil
		}

		_, _ = fmt.Fprintf(out, "%q is not one of: %s\n", answer, strings.Join(options, ", "))
	}
}

func match(in *bufio.Reader, out io.Writer, question string, pattern *regexp.Regexp) (string, error) {
	prompt := question + ": "

	for {
		answer, err := read(in, out, prompt)
		if err != nil {
			return "", err
		}

		if pattern.MatchString(answer) {
			return answer, nil
		}

		_, _ = fmt.Fprintf(out, "%q is not a valid answer\n", answer)
	}
}

func read(in *bufio.Reader, out io.Writer, prompt string) (string, error) {
	_, _ = fmt.Fprint(out, prompt)

	line, err := in.ReadString('\n')
	if err != nil && line == "" {
		return "", ErrNoAnswer
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
