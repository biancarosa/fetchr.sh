# Changelog and Release Management Guide

This guide explains how to use the conventional commit and changelog generation features in fetchr.sh.

## Overview

The project uses:
- **[Conventional Commits](https://www.conventionalcommits.org/)** for structured commit messages
- **[git-cliff](https://git-cliff.org/)** for automated changelog generation
- **[Semantic Versioning](https://semver.org/)** for version numbering
- **Make commands** for automation and consistency

## Quick Start

1. **Install tools**:
   ```bash
   make install-changelog-tools
   ```

2. **Initialize configuration**:
   ```bash
   make changelog-init
   ```

3. **Make a conventional commit**:
   ```bash
   make commit-feat msg="add user authentication"
   ```

4. **Create a release**:
   ```bash
   make release-patch  # for bug fixes
   make release-minor  # for new features
   make release-major  # for breaking changes
   ```

## Conventional Commit Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Supported Types

| Type | Description | Semantic Version Bump |
|------|-------------|----------------------|
| `feat` | New feature | Minor |
| `fix` | Bug fix | Patch |
| `docs` | Documentation only | None |
| `style` | Code style changes (formatting, etc.) | None |
| `refactor` | Code change that neither fixes a bug nor adds a feature | None |
| `test` | Adding or modifying tests | None |
| `chore` | Maintenance tasks, dependency updates | None |
| `perf` | Performance improvements | Patch |
| `security` | Security improvements | Patch |
| `ci` | CI/CD changes | None |
| `build` | Build system changes | None |
| `revert` | Reverting previous commits | Context dependent |

### Breaking Changes

Add `!` after the type or include `BREAKING CHANGE:` in the footer:

```bash
# Breaking change with !
make commit-breaking msg="remove deprecated API endpoints"

# Or with BREAKING CHANGE footer
git commit -m "feat: new API structure

BREAKING CHANGE: API endpoints have been restructured"
```

## Commit Helper Commands

### Basic Commands

```bash
# Feature commits
make commit-feat msg="add user dashboard"
make commit-feat msg="implement file upload" scope="api"

# Bug fixes
make commit-fix msg="resolve memory leak in proxy"
make commit-fix msg="fix CORS headers" scope="dashboard"

# Documentation
make commit-docs msg="update API documentation"

# Refactoring
make commit-refactor msg="extract authentication logic"

# Performance improvements
make commit-perf msg="optimize database queries"

# Security fixes
make commit-security msg="sanitize user inputs"

# Tests
make commit-test msg="add integration tests for proxy"

# Maintenance
make commit-chore msg="update dependencies"

# Breaking changes
make commit-breaking msg="redesign API structure"
make commit-breaking msg="remove deprecated endpoints" scope="api"
```

### Scopes

Use scopes to specify which part of the codebase is affected:

- `api` - Backend API changes
- `dashboard` - Frontend dashboard changes
- `proxy` - Proxy server functionality
- `auth` - Authentication/authorization
- `db` - Database related changes
- `ci` - CI/CD pipeline changes
- `docs` - Documentation changes

## Changelog Generation

### Commands

```bash
# Initialize git-cliff configuration (one-time setup)
make changelog-init

# Generate full changelog from entire git history
make changelog

# Update changelog with commits since last tag
make changelog-update

# Check if recent commits follow conventional format
make check-conventional-commits
```

### Configuration

The changelog is configured in `cliff.toml`. Key features:

- **Emoji grouping**: Features 🚀, Bug Fixes 🐛, etc.
- **Automatic linking**: Issues and mentions become clickable links
- **Semantic grouping**: Related changes are grouped together
- **Breaking change protection**: Important changes aren't skipped

### Customization

Edit `cliff.toml` to customize:

```toml
# Add your repository information
[remote.github]
owner = "yourusername"
repo = "fetchr.sh"

# Customize commit preprocessing for your issue tracker
[[git.commit_preprocessors]]
pattern = '\((\w+\s)?#([0-9]+)\)'
replace = "([#${2}](https://github.com/yourusername/fetchr.sh/issues/${2}))"
```

## Release Management

### Semantic Release Commands

```bash
# Patch release (1.0.0 → 1.0.1) - bug fixes
make release-patch

# Minor release (1.0.0 → 1.1.0) - new features
make release-minor

# Major release (1.0.0 → 2.0.0) - breaking changes
make release-major

# Prerelease (1.0.0 → 1.0.1-alpha.0)
make release-prerelease name=alpha
make release-prerelease name=beta
make release-prerelease name=rc
```

### Manual Release

```bash
# Specify exact version
make release version=1.2.3
```

### Release Process

Each release command automatically:

1. ✅ **Validates** conventional commits since last tag
2. 📝 **Updates** changelog with new entries
3. 🧪 **Runs** all tests and checks (`make check`)
4. 📦 **Commits** changelog changes
5. 🏷️ **Creates** git tag with version
6. 📋 **Provides** next steps for publishing

### After Release

```bash
# Push the release
git push origin main --tags

# Create GitHub/GitLab release from tag
# Build and publish artifacts as needed
```

## Quality Checks

### Pre-commit Validation

```bash
# Check recent commits follow conventional format
make check-conventional-commits

# Run all quality checks
make check
```

### Integration with CI/CD

Add to your workflow:

```yaml
- name: Validate conventional commits
  run: make check-conventional-commits

- name: Generate changelog
  run: make changelog

- name: Create release
  if: github.event_name == 'push' && startsWith(github.ref, 'refs/tags/')
  run: |
    make changelog-update
    # Deploy/publish steps
```

## Workflow Examples

### Feature Development

```bash
# 1. Create feature branch
git checkout -b feature/user-auth

# 2. Make changes and commit
make commit-feat msg="add JWT authentication" scope="api"
make commit-feat msg="add login form" scope="dashboard"
make commit-test msg="add auth integration tests"

# 3. Before merging, check commits
make check-conventional-commits

# 4. Merge to main
git checkout main
git merge feature/user-auth
```

### Bug Fix Release

```bash
# 1. Fix the bug
make commit-fix msg="resolve proxy timeout issue" scope="proxy"

# 2. Create patch release
make release-patch

# 3. Push release
git push origin main --tags
```

### Major Release with Breaking Changes

```bash
# 1. Make breaking changes
make commit-breaking msg="redesign API endpoints" scope="api"
make commit-feat msg="add new dashboard layout" scope="dashboard"

# 2. Update documentation
make commit-docs msg="update API migration guide"

# 3. Create major release
make release-major

# 4. Push and announce
git push origin main --tags
```

## Best Practices

### Commit Messages

- **Be specific**: "fix timeout in proxy" vs "fix bug"
- **Use imperative mood**: "add feature" not "added feature"
- **Include scope when relevant**: `feat(api): add authentication`
- **Reference issues**: "fix memory leak (fixes #123)"

### Releases

- **Patch releases**: Only bug fixes, no new features
- **Minor releases**: New features that don't break existing functionality
- **Major releases**: Breaking changes that require user action
- **Prereleases**: Test versions before stable release

### Changelog

- **Keep it current**: Update changelog regularly, not just at release time
- **Be descriptive**: Users should understand what changed and why
- **Link to issues**: Help users find more context

## Troubleshooting

### Common Issues

**`git-cliff` not found**:
```bash
make install-changelog-tools
```

**Non-conventional commits detected**:
```bash
# Check which commits need fixing
make check-conventional-commits

# Fix commit messages if needed
git rebase -i HEAD~3  # Interactive rebase last 3 commits
```

**Changelog not updating**:
```bash
# Ensure you have conventional commits since last tag
git log --oneline $(git describe --tags --abbrev=0)..HEAD

# Force regenerate full changelog
make changelog
```

**Version calculation incorrect**:
- Check that commits follow conventional format
- Verify breaking changes are marked with `!` or `BREAKING CHANGE:`
- Review `cliff.toml` configuration

### Getting Help

- **Conventional Commits**: https://www.conventionalcommits.org/
- **git-cliff Documentation**: https://git-cliff.org/docs/
- **Semantic Versioning**: https://semver.org/
- **Project Issues**: https://github.com/yourusername/fetchr.sh/issues 