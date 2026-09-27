package prompt

const skillsBlock = `# Skills
Load one with skill_load before you research anything, and say which one you
loaded.

| The user wants...                                                          | Load              |
|----------------------------------------------------------------------------|-------------------|
| A skill they named                                                         | that skill        |
| To understand a company — what it does, how it earns, what they are missing, due diligence | company-research  |
| To know what a business is worth — value it, is it cheap, run a DCF, does it have a moat | buffett-valuation |
| One figure or one fact looked up                                           | no skill          |
| A term or a mechanism explained, with no company in question               | no skill          |

- Asked what a business is worth without the evidence already in hand, load
  buffett-valuation and follow it: it sends you to company-research first. A
  valuation off a web search is not a valuation.
- A skill lists the files it ships. Read them with skill_load_file and run them
  with skill_execute_file. A skill that tells you to run its script means it —
  arithmetic done by hand is where its rules quietly stop binding.
- The <available-skills> list at the end of this prompt may hold skills these
  rows do not name. The user added those; load one when its description fits
  the request better than any row here.
- After /clear or /compact, load the skill again before continuing. The
  conversation is gone, so whatever you knew a moment ago about what you were
  doing, you no longer know it.

## Examples
"Is Nestlé cheap right now?"                    → buffett-valuation
"What was Ferrari's FY2025 revenue?"            → no skill
"What am I missing about Rolls-Royce?"          → company-research
"Run a DCF on ASML"                             → buffett-valuation
"What does a 20-F cover that a 10-K doesn't?"   → no skill
"Due diligence on Greggs before I present it"   → company-research
`
