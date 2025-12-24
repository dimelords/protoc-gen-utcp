# Dependency Check Skill

This skill checks Go module dependencies and validates them against current best practices using Context7.

## Usage

```
/dependency-check
```

Or just mention "check dependencies" in your conversation with Claude Code.

## What it checks

1. **Go module verification** - Ensures go.sum matches go.mod
2. **Outdated dependencies** - Lists direct dependencies with available updates
3. **Context7 validation** - Checks if dependencies follow current best practices from official docs
4. **Vulnerability scanning** - Checks for known security issues (if govulncheck installed)
5. **Module tidiness** - Verifies go.mod doesn't need tidying
6. **Breaking changes** - Warns about major version updates
7. **Deprecated packages** - Checks Context7 for deprecation notices

## Install dependencies (optional)

For enhanced checks:

```bash
# Install govulncheck for vulnerability scanning
go install golang.org/x/vuln/cmd/govulncheck@latest
```

## Automatic checks

This skill is automatically triggered:
- Before commits (if you ask Claude to commit)
- When you modify go.mod or go.sum
- When explicitly invoked with /dependency-check
- When adding new dependencies

## How Context7 integration works

When checking dependencies, the skill:
1. Identifies major dependencies (protobuf, grpc, etc.)
2. Queries Context7 for current recommended versions
3. Compares your versions against official documentation
4. Warns about deprecated APIs or patterns
5. Suggests migrations if breaking changes detected

## Configuration

Add to your `.claude/CLAUDE.md` to customize:

```markdown
## Dependency Management

- Check dependencies before every commit
- Use Context7 to validate against official docs
- Warn on major version updates
- Auto-suggest security fixes
- Check for deprecated APIs
```

## Example output

```
🔍 Checking dependencies...
✅ go.mod and go.sum are in sync
📦 Checking google.golang.org/protobuf v1.36.11
   ✅ Latest stable version (Context7)
   ✅ No known vulnerabilities
⚠️  golang.org/x/tools has update available: v0.15.0 -> v0.16.0
   📚 Checking Context7 for breaking changes...
   ✅ No breaking changes in patch update
✅ All dependency checks passed!
```
