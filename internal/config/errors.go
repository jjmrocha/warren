package config

import "errors"

var (
	ErrConfigNotFound   = errors.New("config not found")
	ErrInvalidProvider  = errors.New("provider is not openrouter, ollama or anthropic")
	ErrInvalidEffort    = errors.New("effort is not off, low, medium or max")
	ErrInvalidMCPName   = errors.New("mcp name is not a bare name")
	ErrInvalidSkillName = errors.New("skill name is not a bare name")
	ErrMissingModel     = errors.New("model is not set")
	ErrMissingAPIKey    = errors.New("api key variable is not set")
	ErrUnknownMCP       = errors.New("mcps-on names a server that is not registered")
)
