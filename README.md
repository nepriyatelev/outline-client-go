# outline-client-go

## Guidelines

- Go testing instructions: `/.github/copilot-go-testing-instructions.md`
- GoDoc comments instructions: `/.github/copilot-go-doc-comments-instructions.md`
- Conventional commit instructions: `/.github/copilot-commit-instructions.md`
- Commit hook: `./githooks/commit-msg`
- Commit template: `./.gitmessage.txt`
- JetBrains live templates: `./.idea/liveTemplates/CopilotTemplates.xml`

Setup (optional):

- Configure Git commit template:
  ```bash
  git config commit.template .gitmessage.txt
  ```
- Use the provided commit-msg hook:
  ```bash
  git config core.hooksPath githooks
  ```
  Make sure the hook is executable (`chmod +x githooks/commit-msg`).
