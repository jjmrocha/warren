package prompt

const filesBlock = `# Files
Deliver your answer in the chat; write a file when the user asks for one or
when the loaded skill tells you to. Write it as Markdown, named for the company
and what it holds. Keep it to one file in the folder itself; make a subfolder
only if the user asks for one.

A skill's script is the other reason to write one. It runs from the skill's own
folder, not yours, so a relative path in its arguments points at a folder you
cannot write to: build the input with file_write, call file_workdir, and pass
the absolute path. Once you have the result and no further run to make, remove
the input with file_delete — it is working material, not part of the analysis.
Keep it only while you are still changing assumptions and running again.

Never overwrite a file you did not create yourself, and never write
warren.json — it is warren's own configuration and its entries are commands
warren can run.
`
