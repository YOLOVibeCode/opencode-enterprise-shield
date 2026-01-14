# Changelog

All notable changes to Enterprise Shield will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2026-01-14

### Added
- Initial release of Enterprise Shield plugin
- Sanitization engine with regex-based pattern matching
- Desanitization engine for response processing
- Compliance detector for PII and secrets (SSN, credit cards, API keys)
- Session management with per-user isolation
- Policy engine with RBAC (unrestricted, sanitized_only, blocked)
- Audit logging with Ed25519 cryptographic signatures
- AES-256-GCM encryption for session data
- Support for multiple LLM providers (OpenAI, Anthropic, Azure, Google)
- Default sanitization rules for:
  - Server and database names
  - IP addresses (RFC 1918 private ranges)
  - Table names
  - Connection strings
  - File paths (Windows and UNC)
  - Internal hostnames
- Compliance detection for:
  - Social Security Numbers (SSN)
  - Credit card numbers (with Luhn validation)
  - AWS access keys
  - GitHub tokens (PAT and OAuth)
  - OpenAI API keys
  - Anthropic API keys
  - Private keys (RSA, EC, DSA)
  - Passwords in code
  - Bearer tokens
  - Slack tokens
  - Azure storage keys
- CLI commands: version, init, scan, process, serve
- Configuration via YAML file
- Installation scripts (install.sh, uninstall.sh)
- Homebrew formula for macOS/Linux
- GitHub Actions workflows for CI/CD
- OpenCode plugin manifest and JavaScript wrapper
- Comprehensive test suite (24 tests)
- Documentation and usage examples

### Security
- All session data encrypted with AES-256-GCM
- Audit logs cryptographically signed with Ed25519
- Fail-secure design (blocks on error)
- No sensitive data in logs (redacted)
- Session isolation between users
- Zero-knowledge architecture

## [0.1.0] - 2026-01-13

### Added
- Initial project structure
- Core type definitions
- Basic sanitization engine prototype

[Unreleased]: https://github.com/YOLOVibeCode/opencode-enterprise-shield/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/YOLOVibeCode/opencode-enterprise-shield/releases/tag/v1.0.0
[0.1.0]: https://github.com/YOLOVibeCode/opencode-enterprise-shield/releases/tag/v0.1.0

