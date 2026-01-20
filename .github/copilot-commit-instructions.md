# Commit Message Instructions for GitHub Copilot

**Purpose**: Generate high-quality, idiomatic Conventional Commits that communicate intent clearly and enable automated tooling (changelog generation, semantic versioning, automated releases).

## 🎯 Copilot Commit Generation Priorities

**ALWAYS follow these in order:**

1. **ALWAYS** write the message in English only
2. **ALWAYS** use a valid Conventional Commit type
3. **ALWAYS** use imperative mood (present tense) in description
4. **ALWAYS** keep description concise (~50-72 characters)
5. **ALWAYS** use lowercase for type and scope
6. **PREFER** including a scope in parentheses (but omit if affects multiple components)
7. **PREFER** including a body when context is helpful (what/why, not how)
8. **INCLUDE** footers for references (Refs, Closes, Fixes) or BREAKING CHANGE details
9. **USE** `!` or `BREAKING CHANGE:` footer if the commit breaks public APIs
10. **AVOID** non-standard types unless your tooling explicitly supports them
11. **AVOID** vague words like "update", "change", "stuff", "misc"
12. **AVOID** Russian or mixed-language text
13. **AVOID** scope if change affects entire project or multiple components
14. **AVOID** periods at the end of description
15. **AVOID** redundant information between description and body
16. **AVOID** emojis or decorative characters in the subject line

Note: Descriptions should start lowercase unless they contain proper nouns or widely-used abbreviations (e.g., JWT, HTTP, Go), which may remain uppercase.

---

## Commit Message Structure

Every commit message follows this structure:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Subject Line (REQUIRED)

Format: `type(scope): description`

**Rules:**
- **Type**: REQUIRED, must be one of the valid types
- **Scope**: OPTIONAL, enclosed in parentheses
- **Description**: REQUIRED, concise summary of changes
- **Total length**: Should fit on one line (~50-72 characters)
- **Case**: lowercase only
- **Ending**: NO period at the end
- **Mood**: MUST use imperative mood (present tense)

**Valid types:**
- `feat` — new feature
- `fix` — bug fix
- `docs` — documentation changes
- `style` — code formatting (no logic change)
- `refactor` — code refactoring (no behavior change)
- `perf` — performance improvement
- `test` — test additions/modifications
- `build` — build system/tooling changes
- `ci` — CI/CD configuration changes
- `chore` — dependency updates, version bumps
- `revert` — reverts a previous commit

---

## Imperative Mood (Present Tense) - CRITICAL

Description MUST use imperative mood, as if commanding someone to apply the commit:

### ❌ WRONG (Past Tense)
```
feat: added data limit support
fix: fixed race condition
refactor: refactored authentication logic
```

### ✅ CORRECT (Imperative/Present)
```
feat: add data limit support
fix: prevent race condition
refactor: simplify authentication logic
```

**Common verbs (imperative):**
- add, remove, delete
- implement, create, initialize
- fix, prevent, resolve, address
- update, modify, change, adjust
- optimize, improve, enhance
- simplify, refactor, reorganize
- support, enable, disable
- document, clarify, explain
- test, verify, validate

---

## Scope Guidelines

### When to Use Scope

Scope clarifies which part of the codebase is affected:

```
✅ feat(api): add JWT authentication
✅ fix(parser): handle edge case in array parsing
✅ docs(readme): add installation instructions
✅ test(auth): add OAuth2 flow tests
```

### When to OMIT Scope

Omit scope if the change affects:
- Multiple components simultaneously
- The entire project/codebase
- General tooling or infrastructure
- Multiple packages/modules

```
✅ fix: prevent race condition across all services
✅ docs: update contributing guidelines
✅ ci: update GitHub Actions workflow
✅ chore: upgrade Go version to 1.21
```

### Scope Format Rules

- **Single word** when possible: `feat(api)`; if necessary, hyphenate: `feat(http-client)`
- **Lowercase only**: `feat(auth)` not `feat(Auth)` or `feat(AUTH)`
- **Hyphenated** if needed: `feat(access-control)` not `feat(accessControl)`
- **No slashes/paths**: `feat(http)` not `feat(src/http)` or `feat(core/http)`
- **No generic names**: avoid `core`, `main`, `utils` — be specific
- **One logical unit**: represent a single module or package

### Common Scope Examples

- **By layer**: api, cli, core, database, frontend, gateway
- **By feature**: auth, payment, cache, logging, metrics, scheduler
- **By module**: parser, tokenizer, serializer, validator, formatter
- **By domain**: user, account, subscription, billing, notification
- **Infrastructure**: build, ci, deps (dependencies), docker
- **Documentation**: readme, docs, changelog, contributing

---

## Description Rules

### Length and Format

```
feat(api): add JWT authentication with refresh tokens
|--------| |--------------------------------------|
   type                    description
                       ~50-72 characters ideal
```

**Rules:**
- Keep it short and specific (~50-72 characters)
- Do NOT end with a period
- Use imperative mood (command form)
- Be specific, not vague

### ❌ Bad Descriptions

```
❌ "update dependencies" (too vague)
❌ "fix stuff" (unclear)
❌ "add feature to support thing" (too generic)
❌ "misc improvements" (meaningless)
❌ "Add data limit support." (period, not imperative; starts with uppercase)
❌ "added support for new API" (past tense)
❌ "feat(api): add JWT authentication 🚀" (emoji not allowed)
```

### ✅ Good Descriptions

```
✅ "add JWT authentication with refresh tokens"
✅ "fix nil pointer in retry logic"
✅ "optimize database query performance by 40%"
✅ "refactor authentication middleware"
✅ "prevent session expiration on token refresh"
✅ "support custom headers in HTTP client"
```

---

## Body (Optional but Recommended)

Use body when additional context is helpful.

### When to Include Body

- Change requires explanation beyond the description
- Complex implementation details
- Motivation for the change
- Impact on other systems
- Non-obvious consequences

### Body Format Rules

- **Separation**: Begin body one blank line after description
- **Line length**: Wrap at 72 characters per line
- **Structure**: Free-form paragraphs, can have multiple paragraphs
- **Content**: Answer "what" and "why", not "how"
- **Multiple paragraphs**: Separate with blank lines

### ❌ Wrong Body Format

```
feat(api): add JWT authentication

Add JWT token support.
Implement refresh token rotation.
Check token expiration.
```

### ✅ Correct Body Format

```
feat(api): add JWT authentication

Introduce JWT token support to replace session-based auth.
This enables stateless authentication across distributed services.

Add support for token refresh to extend sessions without re-authentication.
Implement automatic token rotation for improved security.

Resolves performance issues with session store under high load.
```

### Body Best Practices

- Explain the motivation (why this change was needed)
- Explain the impact (what problem it solves)
- Include relevant context for reviewers
- Do NOT repeat the description
- Do NOT include code snippets (put in commit comments)
- Do NOT explain HOW the code works (use inline comments)

---

## Footers (Optional)

Footers provide metadata and references. Format: `Token: value` or `Token #value`

### Footer Format Rules

- **Separation**: Begin footers one blank line after body (or description if no body)
- **Format**: `Token: value` (colon and space) OR `Token #value` (space and hash)
- **Tokens**: Use hyphens for spaces: `Reviewed-by`, `Co-authored-by`
- **Case**: Tokens are case-insensitive but conventionally Title-Case
- **Multiple**: Separate multiple footers with blank lines or list them consecutively

### Common Footers

| Footer | Format | Purpose | Example |
|--------|--------|---------|---------|
| `Refs` | `Refs: #123` | Reference issue | `Refs: #456` |
| `Closes` | `Closes: #123` | Close issue on merge | `Closes: #789` |
| `Fixes` | `Fixes: #123` | Fix issue on merge | `Fixes: #101` |
| `Reviewed-by` | `Reviewed-by: @name` | Code reviewer | `Reviewed-by: @alice` |
| `Co-authored-by` | `Co-authored-by: Name <email>` | Multiple authors | `Co-authored-by: Bob <bob@example.com>` |
| `BREAKING CHANGE` | `BREAKING CHANGE: description` | Breaking API change | `BREAKING CHANGE: removed deprecated API` |

Inline example using Fixes:

```
fix(parser): handle nil input safely

Avoid panic when input is nil by adding guard clauses.

Fixes: #101
```

---

## BREAKING CHANGE - CRITICAL

Breaking changes MUST be explicitly marked using either `!` notation or `BREAKING CHANGE:` footer.

### When to Mark as BREAKING CHANGE

- API signature changed (function/method removed or modified)
- Return type changed
- Behavior changed in incompatible way
- Required migration for consumers
- Removed deprecated feature

### Method 1: Using `!` Notation (PREFERRED)

Add `!` immediately before the colon:

```
feat(api)!: remove deprecated authentication endpoint

The old /auth/login endpoint has been removed.
Use /auth/signin endpoint instead.
```

```
chore!: drop support for Go 1.18

Go 1.18 is no longer supported. Upgrade to Go 1.19+.
```

### Method 2: Using `BREAKING CHANGE:` Footer

```
feat(types): rename MetricsEnabled to MetricsSettings

BREAKING CHANGE: exported type was renamed; update imports and usage.
```

### Method 3: Using Both (Maximum Clarity)

```
feat(api)!: remove deprecated authentication endpoint

The old /auth/login endpoint has been removed.
Use /auth/signin endpoint instead.

BREAKING CHANGE: /auth/login endpoint removed; migrate to /auth/signin
```

---

## Semantic Versioning Mapping

Commit types map to semantic versioning:

| Commit Type | Version Bump | Trigger |
|------------|-------------|---------|
| `feat` | MINOR | New feature (+0.1.0) |
| `fix` | PATCH | Bug fix (+0.0.1) |
| `BREAKING CHANGE` (any type) | MAJOR | Breaking change (+1.0.0) |
| `docs`, `style`, `refactor`, `perf`, `test`, `chore`, `build`, `ci` | NONE | No version bump |

**Examples:**
- `feat: add new API` → 1.0.0 → 1.1.0
- `fix: resolve bug` → 1.1.0 → 1.1.1
- `feat!: change API` → 1.1.0 → 2.0.0 (BREAKING)

---

## Type-Specific Guidelines

### `feat` — New Feature

Used when adding NEW functionality that users can consume.

```
feat(api): add JWT authentication
feat(cli): add --verbose flag
feat(storage): support S3 backend
```

### `fix` — Bug Fix

Used when fixing bugs, edge cases, or incorrect behavior.

```
fix(parser): handle edge case in array parsing
fix(auth): prevent session expiration on refresh
fix(http): avoid nil pointer in retry logic
```

### `docs` — Documentation

Used for documentation changes (not code).

```
docs: correct spelling of CHANGELOG
docs(readme): add installation instructions
docs(contributing): update contribution guidelines
```

### `style` — Code Formatting

Used for formatting, whitespace, semicolons (NO logic changes).

```
style: format code with gofmt
style(parser): fix indentation
style: remove trailing whitespace
```

### `refactor` — Refactoring

Used when refactoring code (NO behavior changes, NOT an optimization).

```
refactor(auth): simplify token validation
refactor(parser): reorganize functions
refactor(database): extract query builder
```

### `perf` — Performance

Used when optimizing performance (measurable improvement).

```
perf(sorting): optimize comparison function
perf(cache): add in-memory caching layer
perf(database): add indexes to frequently queried columns
```

### `test` — Tests

Used for test additions, modifications, or deletions.

```
test(auth): add OAuth2 flow tests
test: increase coverage to 90%
test(api): add edge case tests for pagination
```

### `build` — Build System

Used for build system, build tools, or build configuration changes.

```
build: add Go module dependencies
build: configure Docker multi-stage build
build: update build script for new targets
```

### `ci` — CI/CD

Used for CI/CD pipeline, workflow, or configuration changes.

```
ci: update GitHub Actions to node 18
ci(gitlab): add security scanning job
ci: parallelize test execution
```

### `chore` — Maintenance

Used for dependency updates, version bumps, or maintenance tasks (no feature or fix).

```
chore: upgrade dependencies
chore: bump Go version to 1.21
chore: update copyright year
```

### `revert` — Revert Previous Commit

Used when reverting a previous commit.

```
revert: add data limit support

This reverts commit abc123.

Refs: #789
```

---

## ❌ Common Mistakes to Avoid

### Mistake 1: Missing Type or Wrong Format

```
❌ "update metrics"
❌ "Update metrics"
❌ ": add metrics"
✅ "feat(metrics): add counter support"
```

### Mistake 2: Description Not Imperative

```
❌ "added JWT authentication"
❌ "adding token refresh"
❌ "will add authentication"
✅ "add JWT authentication"
✅ "implement token refresh"
```

### Mistake 3: Non-English Text

```
❌ "исправление бага в парсере"
❌ "fix(парсер): handle edge case"
✅ "fix(parser): handle edge case"
```

### Mistake 4: Wrong Type Selection

```
❌ "improve: add caching" (not a valid type)
❌ "perf: refactor code" (refactor is not perf)
✅ "perf: add caching layer"
✅ "refactor: simplify code"
```

### Mistake 5: Scope Too Long or Nested

```
❌ "feat(core/client/http): add timeout"
❌ "feat(authentication_and_authorization): add OAuth"
✅ "feat(http): add timeout configuration"
✅ "feat(auth): add OAuth2 support"
```

### Mistake 6: Scope Too Generic

```
❌ "feat(core): add feature"
❌ "fix(utils): fix issue"
✅ "feat(api): add pagination"
✅ "fix(parser): handle null values"
```

### Mistake 7: Blank Line Separation Missing

```
❌ "feat: add feature
Add JWT authentication support.

Refs: #123"

✅ "feat: add JWT authentication

Add support for stateless authentication.

Refs: #123"
```

### Mistake 8: BREAKING CHANGE Not Marked

```
❌ "feat(api): remove deprecated endpoint"
   (No ! or BREAKING CHANGE marker)

✅ "feat(api)!: remove deprecated endpoint"

OR

✅ "feat(api): remove deprecated endpoint

BREAKING CHANGE: /auth/login endpoint removed"
```

### Mistake 9: Description with Period

```
❌ "feat: add JWT authentication."
✅ "feat: add JWT authentication"
```

### Mistake 10: Vague Words

```
❌ "fix: update stuff"
❌ "feat: misc improvements"
❌ "chore: change things"
✅ "fix: prevent nil pointer in parser"
✅ "feat: add data validation"
✅ "chore: upgrade dependencies"
```

---

## 📋 Quick Checklist (Pre-Commit)

Before committing, verify:

- [ ] **English only** — no Russian or mixed languages
- [ ] **Valid type** — one of: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert
- [ ] **Optional scope** — short, relevant, omitted if affects multiple components
- [ ] **Imperative mood** — "add" not "added", "fix" not "fixed"
- [ ] **Concise description** — ~50-72 characters, no period at end
- [ ] **Body present if helpful** — after blank line, answers "why" and "what"
- [ ] **Footers correct** — after blank line, proper format (Token: value)
- [ ] **BREAKING CHANGE marked** — via `!` or footer if applicable
- [ ] **References included** — Refs/Closes for related issues
- [ ] **No vague words** — "update", "stuff", "misc" avoided
- [ ] **Matches specification** — follows Conventional Commits exactly

---

## 📚 Example Messages by Scenario

### Simple Fix

```
fix: prevent race condition in metrics
```

### Feature with Scope

```
feat(api): add pagination with offset and limit
```

### Detailed Feature

```
feat(cache): implement Redis client with connection pooling

Add Redis support to improve performance for frequently accessed data.
Implement automatic connection pooling with configurable pool size.
Support both string and JSON data serialization.

Refs: #234
Closes: #222
```

### Breaking Change with `!`

```
feat(api)!: remove deprecated authentication endpoint

Users must migrate from /auth/login to /auth/signin.
```

### Breaking Change with Footer

```
feat(types): rename UserRole to AccountRole

BREAKING CHANGE: exported type was renamed; update imports and usage.
```

### Multiple Footers

```
fix: prevent nil pointer in retry logic

Add nil-checks around transport and context on retry path.

Reviewed-by: @alice
Co-authored-by: Bob <bob@example.com>
Refs: #456
Closes: #455
```

### Documentation

```
docs: explain cache invalidation strategy
```

### Refactoring

```
refactor(auth): extract token validation into separate function
```

### Performance

```
perf(database): add index to user_id column
```

### Test Addition

```
test(auth): add edge case tests for expired tokens
```

### Dependency Update

```
chore: upgrade Go from 1.20 to 1.21
```

### Revert

```
revert: add data limit support

This reverts commit abc123def456.

Refs: #999
```

---

## 🔗 Useful Resources

- [Conventional Commits Official](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [Git Commit Best Practices](https://git-scm.com/book/en/v2/Git-Basics-Viewing-the-Commit-History)

---

**Last Updated**: January 2026  
**Specification**: Conventional Commits v1.0.0
