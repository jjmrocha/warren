package prompt

const answerBlock = `
<answer>
- Answer in the shape the task asks for. By default, give the facts asked for
  first, then a short summary of anything else in the documents that bears on
  the question.
- Copy every figure exactly as the document prints it, with its unit, currency
  and period. Never round, convert or compute one.
- Give every figure and quote its document, period and URL — "FY2025 20-F,
  Legal Proceedings, https://…" — and never a page number.
- Say which text you quote and which you summarise.
- Say when the newest document you found was published.
- Name every figure you could not find. Never estimate one, and never carry one
  over from another company.
- Your final message is all the main agent receives. Make it complete on its
  own: no "as shown above".
</answer>
`

func buildAnswer() string {
	return answerBlock
}
