# warren

**The AI Financial Analyst**

[![CI](https://github.com/jjmrocha/warren/actions/workflows/ci.yml/badge.svg)](https://github.com/jjmrocha/warren/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/go-1.27%2B-00ADD8)](https://go.dev/dl/)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

warren is a financial analyst you run in your terminal. You start it in a folder, ask it
about a company, and it pulls the market data, reads the filings and the news, and writes up
what the business looks like worth — not where the share price goes next.

Named after [Warren Buffett](https://en.wikipedia.org/wiki/Warren_Buffett).

---

## What it is

**It works in one folder at a time.** warren reads its configuration from `warren.json` in
the directory you start it in, and its file tools are rooted there too. One folder per
company, per portfolio, per piece of research — each with its own model and its own
write-ups. There is nothing in `~/.config`.

**It works through skills.** Asked what a business is worth, warren loads
`buffett-valuation`; asked what a company actually does, it loads `company-research`. The
skills are not warren's own: they come from
[investing-skills](https://github.com/jjmrocha/investing-skills), they live in
`~/.claude/skills`, you install them yourself, and you can edit them.

warren itself is small — the prompt, the configuration and the wiring. The parts under it
are:

| Built on | What it provides |
|---|---|
| [ai-toolkit](https://github.com/jjmrocha/ai-toolkit) | The agent loop, the LLM clients (OpenRouter, Anthropic, Ollama), the tool packs and the skill loader |
| [ai-chat](https://github.com/jjmrocha/ai-chat) | The chat core, the terminal UI and the slash commands |
| [investing-skills](https://github.com/jjmrocha/investing-skills) | `buffett-valuation` and `company-research` — the two skills warren refuses to start without |

You bring the model. warren talks to OpenRouter, Anthropic or a local Ollama, whichever
`warren.json` names.

---

## Install

### 1. Prerequisites

| What | Why | How |
|---|---|---|
| Go 1.27+ | Building warren | [go.dev/dl](https://go.dev/dl/) |
| `donsetch` on `PATH` | Web search, page fetching and crawling. **warren will not start without it** | [donsetch](https://github.com/dondai44423/donsetch) — no key, no account |
| `uvx` on `PATH` | Runs the yfinance MCP server, if you turn it on | [uv](https://github.com/astral-sh/uv) |
| An API key | Unless you run Ollama locally | [OpenRouter](https://openrouter.ai) or [Anthropic](https://console.anthropic.com) |

### 2. Install the skills

warren loads `buffett-valuation` and `company-research` by name at startup and refuses to run
without them. Each must end up as `~/.claude/skills/<name>/SKILL.md`; no other folder is ever
searched, and a missing skill stops warren with its name in the message.

Both live in [investing-skills](https://github.com/jjmrocha/investing-skills). Clone it once
and symlink them, so a `git pull` updates what warren loads:

```bash
git clone https://github.com/jjmrocha/investing-skills.git ~/SOURCES/investing-skills
mkdir -p ~/.claude/skills
ln -sfn ~/SOURCES/investing-skills/buffett-valuation ~/.claude/skills/
ln -sfn ~/SOURCES/investing-skills/company-research  ~/.claude/skills/
```

### 3. Build

```bash
git clone https://github.com/jjmrocha/warren.git && cd warren
make build                      # writes ./bin/warren
```

Copy `./bin/warren` somewhere on your `PATH` — warren is meant to be run from your analysis
folders, not from its own checkout.

`go install github.com/jjmrocha/warren/cmd@latest` works too, but names the binary `cmd`,
after its directory.

### 4. Export your key

```bash
export OPEN_ROUTER_KEY=sk-...
```

Put it in your shell profile so it survives a new terminal. warren never stores the key —
`warren.json` holds the *name* of the variable, and warren reads the variable at startup.

---

## First run

Make a folder for the work and run warren in it:

```bash
mkdir ~/research/acme && cd ~/research/acme
warren
```

It finds no `warren.json`, so it says what it is about to do and asks two questions:

```
warren — The AI Financial Analyst.

First run in this folder: two questions, saved to

  /Users/you/research/acme/warren.json

which you can edit later. warren also loads buffett-valuation and company-research from

  /Users/you/.claude/skills

and will not start without them — github.com/jjmrocha/investing-skills

Provider [anthropic, ollama, openrouter]: openrouter
Model: z-ai/glm-5.3-flash
```

Both paths are the real ones — the folder you started warren in, and your home directory.

From the answers it writes `warren.json` in that folder, with the variable holding your key
filled in from the provider — `OPEN_ROUTER_KEY` for OpenRouter, `ANTHROPIC_API_KEY` for
Anthropic, none for Ollama — and then starts.

Every other folder you run warren in asks the same two questions once, and keeps its own
answers.

---

## Using warren

Run warren in the folder the work belongs in. Its file tools are rooted there and cannot
leave it, so the write-ups it saves land beside the configuration that produced them.

| Command | What it does |
|---|---|
| `/help` | List the commands |
| `/model [name]` | Show the current model, or switch to another from `llm.models` |
| `/effort [level]` | Show or change reasoning effort |
| `/clear` | Reset the conversation |
| `/compact` | Compact the context now, instead of waiting for warren to do it |
| `/mcp [on\|off] [name]` | Show the MCP servers, or start and stop one |
| `/skills` | List the skills warren loaded |
| `/exit` | Quit |

### The tools it always has

These are built into warren. None of them come from `warren.json`, and none can be turned
off.

| Tools | What they do |
|---|---|
| `file_read`, `file_write`, `file_edit`, `file_list`, `file_search`, `file_delete`, `file_workdir` | Read, change and search the files of the folder warren runs in — and nothing outside it |
| `donsetch__web_search`, `donsetch__web_fetch`, `donsetch__web_crawl` | Search the web, read a page, walk a site — served by `donsetch` |
| `current_date`, `current_time`, `time_zone` | Date a report's pull date or a valuation's as-of date, instead of guessing from what the model was trained on |

### Optional — market data

`warren.json` arrives with [yfmcp](https://pypi.org/project/yfmcp/) registered as
`yfinance-mcp`: quotes, fundamentals, statements and filings from Yahoo Finance. It is not
started at launch. Start it in the session with `/mcp on yfinance-mcp`, and `uvx` fetches it
the first time.

Without it warren still works. The filings, the interim reports, the results calls, the
ownership disclosures and the reviews are all on the web, and that is where the research
skill sends it anyway — the server is a shortcut to quotes and financial series, not the
source. A server warren has no tools for is never a reason to stop: the prompt says so, and
the filings are where the answer was coming from anyway.

---

## Configure

warren starts from `warren.json` in the directory you run it in. Not a parent directory, not
`~/.config` — that one file, that one folder. Delete it and the next run asks the two
questions again; edit it and warren starts from what you wrote.

**The file is complete.** warren runs exactly what it says: no merging, no hidden defaults.

```json
{
  "llm": {
    "provider": "openrouter",
    "api-key-env": "OPEN_ROUTER_KEY",
    "model": "z-ai/glm-5.3-flash",
    "models": ["z-ai/glm-5.3-flash", "deepseek/deepseek-v4-pro"],
    "effort": "medium"
  },
  "skills": ["removing-ai-tells"],
  "mcps": {
    "yfinance-mcp": { "command": "uvx", "args": ["yfmcp@latest"], "timeout": 60 }
  },
  "mcps-on": []
}
```

| Key | What it does |
|---|---|
| `llm.provider` | `openrouter`, `ollama` or `anthropic` |
| `llm.base-url` | Overrides the provider's endpoint; omit it to use the standard one |
| `llm.api-key-env` | The **name** of the variable holding the key, never the key itself. Required except on Ollama |
| `llm.model` | The model warren starts with |
| `llm.models` | The models `/model` switches between |
| `llm.effort` | `off`, `low`, `medium` or `max` — how much the model reasons before answering |
| `skills` | Extra skills by name, loaded from `~/.claude/skills` beside warren's own two |
| `mcps` | MCP servers warren registers: `command`, `args`, `env` (variables inherited from warren), `timeout` in seconds — `0` or absent means no limit. The key is a bare name: letters, digits, `_` and `-` |
| `mcps-on` | The servers started at launch. Written empty: the rest stay stopped until you run `/mcp on <name>` |

A file that names an unknown provider or effort, a skill or an MCP server that is not a bare
name, or an `mcps-on` server that is not in `mcps`, or that leaves the model empty or names an
API-key variable that is not set, stops warren before the session opens — and the message
lists every fault in the file, not just the first. An unknown key in the file is an error too,
so a typo in a section name is caught rather than ignored.

### A `warren.json` you did not write

`mcps` entries are commands warren can run. A `warren.json` that came with a repository, a
shared folder or a colleague can name any command there — so read it before running warren in
that folder, the way you would read a `Makefile`.

`mcps-on` is what makes that a decision rather than an accident: warren starts only the
servers listed there, and the file it writes for you lists none. A server appears in `/mcp`
as `off` until you or the model asks for it.

Your key never sits in the file, so a `warren.json` is safe to commit alongside the analysis
it belongs to.

---

## Troubleshooting

warren validates what it can before the session opens, and the message names the fault.

| Message | Cause | Fix |
|---|---|---|
| `error starting mcp donsetch: exec: "donsetch": executable file not found in $PATH` | The web tools have no server | Install [donsetch](https://github.com/dondai44423/donsetch) and put it on `PATH` |
| `skill folder not found: …` | A skill is missing from `~/.claude/skills` | Install it from [investing-skills](https://github.com/jjmrocha/investing-skills) — every missing one is listed at once |
| `api key variable is not set` | `api-key-env` names a variable with no value | `export` it, or point `api-key-env` at the one you use |
| `config not found: warren.json` | The file vanished between setup and startup | Run warren again; it asks the two questions |
| `no answer to read` | Setup ran with nothing on stdin — a pipe, a redirect, or Ctrl-D at a question | Run warren from a terminal and answer the questions; nothing is left broken, the next run simply asks again |
| `provider is not openrouter, ollama or anthropic` | Unknown `llm.provider` | Use one of the three |
| `mcps-on names a server that is not registered` | A name in `mcps-on` with no entry in `mcps` | Add the server, or drop the name |
| `mcp name is not a bare name` | An `mcps` key with a character outside letters, digits, `_` and `-` | Rename the server |
| `json: unknown field …` | A misspelled key in `warren.json` | Fix the spelling — warren does not ignore keys it does not know |

An `mcps-on` server that fails to start does *not* stop warren. The failure prints before the
TUI opens and `/mcp` shows the server as `off`; `/mcp on <name>` retries it. `donsetch` is the
exception: it is not an `mcps` entry, and warren stops when it is missing.

To start over in a folder, delete its `warren.json` and run warren again.

---

## Development

`make` on its own lists every target:

| Target | Effect |
|---|---|
| `build` | Build warren into `./bin` |
| `clean` | Remove `./bin` |
| `test` | Run all tests |
| `bench` | Run benchmarks |
| `lint` | Run golangci-lint |
| `deps` | Update dependencies |
| `tidy` | Tidy `go.mod` |

CI runs `go test -race ./...`, `golangci-lint` and `govulncheck` on every push and pull
request ([`.github/workflows/ci.yml`](.github/workflows/ci.yml)). Run the `-race` form before
pushing — `make test` does not. golangci-lint must be v2.13.2 or newer; earlier releases are
built with an older Go and refuse this module.

## License

[MIT](LICENSE) © 2026 Joaquim Rocha
