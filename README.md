# warren

**The AI Financial Analyst**

[![CI](https://github.com/jjmrocha/warren/actions/workflows/ci.yml/badge.svg)](https://github.com/jjmrocha/warren/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/go-1.27%2B-00ADD8)](https://go.dev/dl/)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

warren is a financial analyst that runs in your terminal. Name a company and it reads the
filings, the results calls, the ownership disclosures and the reviews, then hands you the
evidence: what the business does, how it makes money, who owns it, what could hurt it. Ask
what the business is worth and it values it the way an owner would, and shows every number
it used.

It never tells you to buy or sell. You get the evidence and the arithmetic; the view is
yours.

Named after [Warren Buffett](https://en.wikipedia.org/wiki/Warren_Buffett).

---

## What you can ask it

```
What am I missing about Rolls-Royce?
Due diligence on Greggs before I present it
Is Nestlé cheap right now?
Run a DCF on ASML
What was Ferrari's FY2025 revenue?
What does a 20-F cover that a 10-K doesn't?
```

warren works out which of four kinds of request you made and answers in the matching way.

| You want... | warren... | You get |
|---|---|---|
| **To understand a company**: what it does, how it earns, what you're missing, due diligence | runs *company-research* | A full research report |
| **To know what a business is worth**: value it, is it cheap, run a DCF, does it have a moat, can you trust management | runs *buffett-valuation*, doing the research first if it hasn't already | Four numbers, each with its reasoning |
| **One figure or one fact** | looks it up | The figure, with the document and period it came from |
| **A term or a mechanism explained** | answers | The explanation |

You can also start either job directly with `/company-research <company>` or
`/buffett-valuation <company>`.

---

## Research: the whole picture, sourced

Ask warren to research a company and it produces one complete report, whether or not you
already know the company. That is deliberate. Every report covers the same ground, so for a
company you know well, the sections that surprise you are the ones you were missing.

**What it reads, every time:**

- the latest annual report, **footnotes included** (10-K, 20-F, or the local equivalent)
- the latest interim report
- the last two results calls, prepared remarks and Q&A
- ownership disclosures
- a year of material events and news
- employee reviews (Glassdoor and similar)
- customer and product reviews

The footnotes matter most. The narrative sections are written by investor relations; the
footnotes are written by lawyers and accountants. That is where you find that three
distributors are half of revenue, that a litigation reserve moved and why, or that an
accounting change makes this year incomparable with last.

**What the report covers:** a factual TL;DR, then what the company does, products and
services, customers and their concentration, geography, history, ownership, recent
developments, how revenue is actually earned, evidence bearing on durability, strategy versus
where the capital actually went, risks, and what employees and customers say.

**It ends with *What I could not establish*.** A line warren could not source never appears
in the body. There are no blanks and no estimates. The line moves to this closing list, along
with any source it could not reach. A long list is itself a finding: that company discloses
little.

---

## Valuation: four numbers you can check

Ask what a business is worth and warren returns exactly this:

| | |
|---|---|
| **Header** | Enterprise value (market cap + debt − cash), its as-of date, and where the price came from |
| **1. Intrinsic value** | One number for the whole company, from ten years of discounted owner earnings plus a terminal value. Shows the bridge from reported net income, how maintenance capex was derived, the growth rates and their source, the discount rate, and the terminal multiple |
| **2. Intrinsic value × 70%** | The same number with a margin of safety applied |
| **3. Moat tier** | Inevitable, Formidable, Strong or None. Names the mechanism behind it and what would break it |
| **4. Management** | Green flags met and red flags present, each tied to a disclosure |

You compare enterprise value with intrinsic value yourself. warren puts the two numbers side
by side and stops.

**The rules are conservative, and they are fixed:**

- Growth is capped at 15% for years 1–5 and 10% for years 6–10, whatever the company has
  done recently.
- The discount rate never falls below 10%, so cheap money cannot inflate the result.
- Stock-based compensation counts as a cost, not as a non-cash charge to add back.
- Maintenance capex is derived from the filings rather than assumed equal to total capex.
- The terminal value is an exit multiple set by moat tier, from 15× down to 9×.

The arithmetic runs in a script that ships with the skill, not in the model's head, so these
rules cannot be argued away by a persuasive growth story. Every business gets a number. Where
an input is weak, warren says so on the line where it is used, rather than refusing or
giving you a range.

---

## How warren works

- **Filings over headlines.** It finds the company's filing home first (EDGAR for a US
  filer; otherwise the investor relations site, the national regulator or the company
  registry). It prefers the filing to the press release, and the press release to the news.
- **Nothing from memory.** Every figure comes with its document and period, such as "FY2025
  20-F, Legal Proceedings". It cites no page numbers, because those are the easiest to get
  wrong.
- **Dated, not "current".** It reads today's date from the clock, dates every figure by the
  period its source covers, and tells you how recent its newest document is.
- **Missing means missing.** A figure it could not find is named as missing. It is never
  estimated or borrowed from a similar company.
- **Ambiguity is a question.** A name or ticker that fits more than one company gets a
  question back, not a guess.
- **It tells you its sources.** Every answer says which sources it reached and which it
  could not.

Answers arrive in the chat. warren writes a file only when you ask for one, or when a skill
needs a working file for its script. Files are Markdown, named for the company and their
contents, and saved in the folder you started warren in.

---

## Getting started

### 1. What you need

| What | Why |
|---|---|
| An API key for [OpenRouter](https://openrouter.ai) or [Anthropic](https://console.anthropic.com), or a local [Ollama](https://ollama.com) | The model warren thinks with |
| [`donsetch`](https://github.com/dondai44423/donsetch) on your `PATH` | Web search, page reading and crawling. No key or account needed. **warren will not start without it** |
| [`uv`](https://github.com/astral-sh/uv) on your `PATH` | Runs the valuation script, and the optional market-data server |
| [Go 1.27+](https://go.dev/dl/) | To build warren |

### 2. Install the two skills

warren's research and valuation methods live in
[investing-skills](https://github.com/jjmrocha/investing-skills), as plain files you can read
and edit. warren loads them from `~/.claude/skills` and refuses to start without them. Clone
them once and link them in, so that a `git pull` updates what warren uses:

```bash
git clone https://github.com/jjmrocha/investing-skills.git ~/SOURCES/investing-skills
mkdir -p ~/.claude/skills
ln -sfn ~/SOURCES/investing-skills/buffett-valuation ~/.claude/skills/
ln -sfn ~/SOURCES/investing-skills/company-research  ~/.claude/skills/
```

### 3. Build warren

```bash
git clone https://github.com/jjmrocha/warren.git && cd warren
make build                      # writes ./bin/warren
```

Copy `./bin/warren` to a folder on your `PATH`. You run warren from your research folders,
not from its source checkout.

### 4. Export your key

```bash
export OPEN_ROUTER_KEY=sk-...   # or ANTHROPIC_API_KEY for Anthropic
```

Add it to your shell profile so it is there in every new terminal. warren never stores the
key. Its configuration holds only the *name* of the variable.

### 5. First run

Make a folder for the work and start warren in it:

```bash
mkdir ~/research/acme && cd ~/research/acme
warren
```

The first time, warren asks two questions:

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

It saves your answers to `warren.json` in that folder and opens the chat. Any other folder
asks the same two questions once and keeps its own answers.

---

## Using warren

**One folder per piece of work.** Use a folder per company, per portfolio or per question.
Each folder keeps its own configuration, its own write-ups and its own saved sessions.
warren's file tools can only reach the folder you started it in.

### Commands

| Command | What it does |
|---|---|
| `/company-research [request]` | Research a company before forming a view |
| `/buffett-valuation [request]` | Value a business as a whole |
| `/model [name]` | Show the current model, or switch to another listed in `warren.json` |
| `/effort [level]` | Show or change how much the model reasons before answering |
| `/mcp [on\|off] [name]` | Show the data servers, or start and stop one |
| `/skills` | List the skills warren loaded |
| `/export` | Save the conversation so you can resume it later |
| `/compact` | Shrink the conversation now, so a long session keeps going |
| `/clear` | Start a fresh conversation |
| `/help` | List the commands |
| `/exit` | Quit |

After `/clear` or `/compact`, warren loads its skill again before it continues.

### Saving and resuming a session

`/export` saves the conversation to `sessions/<id>.json` in the current folder and prints
the id. Exporting again overwrites the same file. `/clear` starts a new session, so its next
export goes to a new file. Only you can read the files, but they hold the whole conversation,
so keep that in mind before you share the folder.

To pick a session up again, start warren in the same folder with its id:

```sh
warren -resume <id>
```

The conversation is replayed on screen and the model continues from it.

### Market data (optional)

Every new `warren.json` includes [yfmcp](https://pypi.org/project/yfmcp/), registered as
`yfinance-mcp`. It provides quotes, fundamentals, statements and filings from Yahoo Finance.
It is off by default. Start it in a session with `/mcp on yfinance-mcp`; `uv` fetches it the
first time.

warren works fully without it. Filings, results calls and disclosures are all on the web, and
that is where the research goes anyway. The server is a shortcut to price data and financial
series, not a source warren depends on.

### Your own skills

List more skills from `~/.claude/skills` under `skills` in `warren.json`, for example
`"skills": ["removing-ai-tells"]`. warren uses one whenever its description fits your request
better than research or valuation. `/skills` shows what is loaded.

---

## Configure

Each folder's `warren.json` is the whole configuration. There is no merging and no hidden
default. Delete it and the next run asks the two questions again.

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
| `llm.base-url` | Overrides the provider's endpoint. Leave it out to use the standard one |
| `llm.api-key-env` | The **name** of the variable holding your key, never the key itself. Required except with Ollama |
| `llm.model` | The model warren starts with |
| `llm.models` | The models `/model` can switch between |
| `llm.effort` | `off`, `low`, `medium` or `max`: how much the model reasons before answering |
| `skills` | Extra skills by name, loaded from `~/.claude/skills` |
| `mcps` | Data servers warren can start: `command`, `args`, `env` (variables passed through from your shell), `timeout` in seconds (`0` or absent means no limit). Names use letters, digits, `_` and `-` |
| `mcps-on` | The servers started at launch. Empty by default: the rest wait for `/mcp on <name>` |

warren checks the file before it opens a session. A misspelled key is an error, not
something silently ignored, and a bad value stops warren with every faulty value listed at
once.

### Safety

The entries under `mcps` are commands that warren runs. Before you run warren in a folder
whose `warren.json` came from someone else (a shared drive, a repository, a colleague),
read that file the way you would read a script. Only the servers listed in `mcps-on` start on
their own, and the file warren writes for you lists none.

warren itself cannot change that file. Its file tools can read `warren.json` and `sessions/`,
but they refuse to write, edit or delete either. A web page that tries to talk the model into
changing what warren runs gets nowhere.

Your key is never written to `warren.json`, so the file is safe to keep alongside the
research it belongs to.

---

## Troubleshooting

| Message | Cause | Fix |
|---|---|---|
| `error starting mcp donsetch: exec: "donsetch": executable file not found in $PATH` | The web tools have no server | Install [donsetch](https://github.com/dondai44423/donsetch) and put it on your `PATH` |
| `skill folder not found: …` | A skill is missing from `~/.claude/skills` | Install it (step 2 above). Every missing skill is listed at once |
| `api key variable is not set` | `api-key-env` names a variable with no value | `export` it, or point `api-key-env` at the one you use |
| `no answer to read` | Setup ran with nothing to read, for example through a pipe, a redirect, or Ctrl-D at a question | Run warren in a terminal and answer the questions. Nothing was written; the next run asks again |
| `provider is not openrouter, ollama or anthropic` | Unknown `llm.provider` | Use one of the three |
| `mcps-on names a server that is not registered` | A name in `mcps-on` has no entry in `mcps` | Add the server, or remove the name |
| `mcp name is not a bare name` | An `mcps` name uses something other than letters, digits, `_` and `-` | Rename the server |
| `json: unknown field …` | A misspelled key in `warren.json` | Fix the spelling |
| `resume: session not found` | No `sessions/<id>.json` in this folder | Start warren in the folder the session was exported from |
| `resume: invalid session id` | The id is not one `/export` printed | Copy the id exactly as `/export` showed it |
| `config not found: warren.json` | The file was removed between setup and startup | Run warren again; it asks the two questions |

A server in `mcps-on` that fails to start does *not* stop warren. The error is printed
before the chat opens, `/mcp` shows the server as `off`, and `/mcp on <name>` tries again.

---

## Development

warren is small: the prompt, the configuration and the wiring. The agent loop, model clients
and tools come from [ai-toolkit](https://github.com/jjmrocha/ai-toolkit); the chat, terminal
UI and slash commands from [ai-chat](https://github.com/jjmrocha/ai-chat).

`make` lists the targets. `make test` runs the tests with the race detector and `make lint`
runs golangci-lint (v2.13.2 or newer). CI runs both, plus `govulncheck`, on every push
([`.github/workflows/ci.yml`](.github/workflows/ci.yml)).

## License

[MIT](LICENSE) © 2026 Joaquim Rocha
