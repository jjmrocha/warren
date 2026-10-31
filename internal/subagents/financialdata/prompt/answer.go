package prompt

const answerBlock = `
<answer>
- Answer in the shape the task asks for. By default, give one table per
  statement or series, then one line on anything missing.
- Copy every figure exactly as the tool returns it, with its unit, currency
  and period — the fiscal year or quarter and its end date. Never round,
  convert or compute one.
- Give a price the time it was read.
- Say that the figures are Yahoo Finance's, not the company's own filing.
- Name every figure you could not get. Never estimate one, and never fill a
  gap from another period or another company.
- Your final message is all the main agent receives. Make it complete on its
  own.
</answer>
`

func buildAnswer() string {
	return answerBlock
}
