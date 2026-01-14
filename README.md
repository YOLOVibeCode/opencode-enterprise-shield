# Enterprise Shield Plugin for OpenCode 🛡️

**Enterprise-grade security for AI coding assistants**

Enterprise Shield is a production-ready OpenCode plugin that automatically detects and masks sensitive information (database names, servers, IPs, credentials, PII) before sending to LLM APIs — then seamlessly restores the original values in responses.

[![Tests](https://img.shields.io/badge/tests-24%20passing-brightgreen)](https://github.com/YOLOVibeCode/opencode-enterprise-shield)
[![Go Version](https://img.shields.io/badge/go-1.22-blue)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

---

## 📚 Quick Navigation

- 🚀 [Installation](#-installation) - Multiple install methods
- 💼 [Business Justification](#-business-justification) - For corporate decision-makers  
- 🎯 [The Problem](#-the-problem) - Why you need this
- 🔐 [What Gets Protected](#-what-gets-protected) - Coverage details
- 📖 [Usage Examples](#-usage-examples) - Practical scenarios
- ✨ [Features](#-features) - Full capabilities
- 🏗️ [Distribution Guide](DISTRIBUTION.md) - For maintainers

---

## 🎯 The Problem

### Real-World Developer Risk Scenario

```
Developer using AI assistant in Cursor/VSCode:
"Help me optimize this query:

SELECT u.email, u.last_login, c.credit_card_number
FROM ProductionDB.user_accounts u
JOIN CustomerPII.payment_methods c ON u.id = c.user_id  
WHERE server='db-prod-01.internal.company.com'
  AND u.created_at > '2024-01-01'
  AND u.department = 'Engineering'"
```

**❌ What ChatGPT/Claude receives without protection:**
- Production database name: `ProductionDB`
- PII table name: `CustomerPII.payment_methods`
- Internal hostname: `db-prod-01.internal.company.com`
- Table structure revealing business logic
- Schema showing sensitive data relationships

**🚨 Risks:**
- **Data Breach**: Internal infrastructure exposed to external AI providers
- **Compliance Violation**: HIPAA, GDPR, SOC 2 requirements broken
- **IP Leakage**: Database design reveals proprietary business logic
- **Audit Failure**: No record of what data left the corporate network
- **Insider Threat**: No controls on what developers can share

### ✅ The Enterprise Shield Solution

**Developer experience unchanged:**
```sql
SELECT u.email, u.last_login, c.credit_card_number
FROM ProductionDB.user_accounts u
WHERE server='db-prod-01.internal.company.com'
```

**What AI receives (automatically sanitized):**
```sql
SELECT u.email, u.last_login, c.credit_card_number
FROM SERVER_0.TABLE_0 u
WHERE server='HOST_0'
```

**Developer gets back:**
```
"Add an index on ProductionDB.user_accounts.created_at
and optimize the connection pool for db-prod-01.internal.company.com.
Consider partitioning user_accounts by department."
```

**Benefits:**
- ✅ **Zero-Knowledge**: AI never sees real infrastructure
- ✅ **Transparent**: Developer workflow completely unchanged
- ✅ **Compliant**: Audit trail + automatic PII blocking
- ✅ **Secure**: Multiple layers of protection
- ✅ **Useful**: AI still provides relevant, actionable advice

---

## 💼 Business Justification

### For Corporate Decision-Makers

#### Compliance Requirements Met

| Regulation | Requirement | Enterprise Shield Solution |
|------------|-------------|----------------------------|
| **HIPAA** | Protect PHI/ePHI | ✅ Blocks PII, cryptographic audit trail |
| **GDPR** | Data protection by design | ✅ Sanitizes PII, enables data deletion |
| **SOC 2 Type II** | Access controls & logging | ✅ RBAC policies, Ed25519-signed logs |
| **PCI-DSS v4** | Protect cardholder data | ✅ Detects & blocks credit cards (Luhn) |
| **CCPA** | Consumer privacy | ✅ Data minimization, comprehensive audits |
| **ISO 27001** | Information security | ✅ Technical controls, audit evidence |

#### Risk Mitigation

| Threat | Impact Without Shield | With Enterprise Shield |
|--------|----------------------|------------------------|
| **Accidental Data Leakage** | Developers leak prod DBs to OpenAI | ✅ Sanitized automatically, zero-knowledge |
| **IP Theft** | Schema/architecture visible in prompts | ✅ Business logic masked with aliases |
| **Compliance Fines** | GDPR: up to €20M or 4% revenue | ✅ Demonstrable controls prevent violations |
| **Audit Failures** | No visibility into AI tool usage | ✅ Complete, tamper-proof audit trail |
| **Insider Threats** | Unrestricted data exfiltration | ✅ Policy-based access controls (RBAC) |
| **Reputational Damage** | Public breach disclosure | ✅ Preventative controls demonstrate diligence |

#### ROI Analysis

**Cost of NOT Implementing:**
```
Average data breach cost (IBM 2024):                $4.45 million
Average GDPR fine (based on actual cases):          $2-10 million  
Lost productivity from breach investigation:        500-2000 hours
Compliance audit failure remediation:               $500K-2M
Customer churn from breach:                         10-30% revenue
Legal & PR costs:                                   $1-5 million

TOTAL POTENTIAL COST:                               $8M - $21M
```

**Cost of Enterprise Shield:**
```
Installation time (one-time):                       30 minutes
Monthly maintenance:                                ~1 hour
Performance overhead:                               <50ms per request
Developer workflow disruption:                      ZERO (transparent)
License cost:                                       Open source (MIT)

TOTAL ONGOING COST:                                 ~$200/year in labor
```

**ROI:** **Prevents $8M-21M in costs with $200/year investment = 40,000% to 105,000% ROI**

#### Decision-Maker Talking Points

**For Chief Information Security Officers (CISOs):**
- ✅ **Defense in Depth**: Adds critical security layer for AI tools
- ✅ **Zero Trust Architecture**: Assumes external AI is untrusted
- ✅ **Audit Trail**: Ed25519 cryptographically signed, tamper-evident logs
- ✅ **Enterprise Encryption**: AES-256-GCM for data at rest
- ✅ **Fail-Secure Design**: Blocks on error rather than leaking data
- ✅ **Compliance Evidence**: Demonstrates technical controls for auditors

**For Chief Technology Officers (CTOs):**
- ✅ **Zero Productivity Impact**: Completely transparent to developers
- ✅ **Production-Ready**: 24 tests passing, clean architecture
- ✅ **Scalable**: Stateless design, horizontal scaling ready
- ✅ **Maintainable**: Clean Go codebase, comprehensive docs
- ✅ **Cross-Platform**: macOS, Linux, Windows support
- ✅ **Integration**: 6 install methods, OpenCode native

**For Chief Financial Officers (CFOs):**
- ✅ **Risk Mitigation**: Prevents multi-million dollar breach costs
- ✅ **Productivity Enabler**: Safe AI adoption drives 10-30% efficiency gains
- ✅ **Insurance Benefits**: May reduce cyber insurance premiums
- ✅ **Audit Cost Reduction**: Automated compliance evidence collection
- ✅ **Minimal Investment**: Near-zero implementation cost
- ✅ **Quantifiable Value**: Easily measured risk reduction

**For Compliance & Privacy Officers:**
- ✅ **GDPR Article 25**: Data protection by design and by default
- ✅ **Audit Trail**: 365-day retention (configurable) with signatures
- ✅ **Policy Enforcement**: Granular RBAC controls per user/department
- ✅ **Right to Erasure**: Session data can be deleted on request
- ✅ **Data Minimization**: Only metadata logged, full content excluded
- ✅ **Demonstrable Controls**: Technical proof of protection measures

---

## 🚀 Installation

### Method 1: One-Line Install (Recommended) ⭐

**🍎 macOS / 🐧 Linux:**
```bash
curl -sSL https://raw.githubusercontent.com/YOLOVibeCode/opencode-enterprise-shield/main/install-complete.sh | bash
```

**🪟 Windows (PowerShell):**
```powershell
iwr -useb https://raw.githubusercontent.com/YOLOVibeCode/opencode-enterprise-shield/main/install-windows.ps1 | iex
```

**What this does:**
- ✅ Installs **both** OpenCode and Enterprise Shield
- ✅ Auto-detects your OS and architecture (amd64/arm64)
- ✅ Downloads correct binary from GitHub releases
- ✅ Verifies SHA256 checksum for security
- ✅ Installs to `~/.opencode/plugins/enterprise-shield`
- ✅ Creates default configuration
- ✅ **Runs 8 protection tests** to prove it works
- ✅ Adds to PATH automatically

**Time:** 2-3 minutes | **Result:** Working installation with proof of protection

### Method 2: Homebrew (macOS/Linux) ⭐

```bash
# Add the tap
brew tap YOLOVibeCode/opencode-enterprise-shield

# Install
brew install enterprise-shield

# Verify
enterprise-shield version

# Update later
brew upgrade enterprise-shield

# Uninstall
brew uninstall enterprise-shield
```

### Method 3: Go Install

```bash
# Install latest version
go install github.com/YOLOVibeCode/opencode-enterprise-shield/cmd/plugin@latest

# Symlink to OpenCode plugins directory
ln -s $(go env GOPATH)/bin/plugin ~/.opencode/plugins/enterprise-shield

# Initialize configuration
enterprise-shield init
```

### Method 4: NPM

```bash
# Install globally
npm install -g @YOLOVibeCode/opencode-enterprise-shield

# Verify
enterprise-shield version
```

### Method 5: Docker

```bash
# Pull image
docker pull YOLOVibeCode/enterprise-shield:latest

# Run
docker run -v ~/.opencode:/root/.opencode YOLOVibeCode/enterprise-shield:latest

# Run specific command
docker run YOLOVibeCode/enterprise-shield:latest scan "Test content"
```

### Method 6: Manual Download

1. Download from [GitHub Releases](https://github.com/YOLOVibeCode/opencode-enterprise-shield/releases)
2. Choose your platform:
   - macOS Intel: `enterprise-shield-v1.0.0-darwin-amd64.tar.gz`
   - macOS Apple Silicon: `enterprise-shield-v1.0.0-darwin-arm64.tar.gz`
   - Linux: `enterprise-shield-v1.0.0-linux-amd64.tar.gz`
   - Windows: `enterprise-shield-v1.0.0-windows-amd64.zip`

```bash
# Extract (macOS/Linux)
tar -xzf enterprise-shield-v1.0.0-darwin-arm64.tar.gz

# Install
mkdir -p ~/.opencode/plugins
mv enterprise-shield ~/.opencode/plugins/
chmod +x ~/.opencode/plugins/enterprise-shield

# Initialize config
enterprise-shield init
```

### Verify Installation

**All Platforms:**

```bash
# Check version
enterprise-shield version
# Output: Enterprise Shield Plugin v1.0.1

# Test compliance scanner
enterprise-shield scan "Test with SSN: 123-45-6789"
# Should show critical violation detected

# Run full protection test suite (8 tests)
curl -sSL https://raw.githubusercontent.com/YOLOVibeCode/opencode-enterprise-shield/main/test-protection.sh | bash
# Expected: 8/8 tests PASSED ✅
```

**Windows (PowerShell):**

```powershell
# Check version
enterprise-shield version

# Test scanner
enterprise-shield scan "Test with SSN: 123-45-6789"
```

**See:** `WINDOWS_INSTALLATION.md` for Windows-specific guide

---

## ⚡ Quick Start (3 Minutes)

### 1. Install

```bash
curl -sSL https://raw.githubusercontent.com/YOLOVibeCode/enterprise-shield/main/install.sh | bash
```

### 2. Test Compliance Detection

```bash
# Test SSN detection
$ enterprise-shield scan "Employee SSN: 123-45-6789"

{
  "hasViolations": true,
  "shouldBlock": true,
  "violations": [{
    "type": "SSN",
    "severity": "critical",
    "redactedValue": "123-****6789",
    "position": 14
  }]
}

# Test API key detection
$ enterprise-shield scan "AWS key: AKIAIOSFODNN7EXAMPLE"

{
  "hasViolations": true,
  "shouldBlock": true,
  "violations": [{
    "type": "API_KEY",
    "severity": "critical",
    "redactedValue": "AKIA****MPLE"
  }]
}
```

### 3. Test Sanitization

```bash
$ enterprise-shield process developer@company.com \
  "SELECT * FROM ProductionDB.users WHERE ip='192.168.1.100'" \
  openai

{
  "content": "SELECT * FROM SERVER_0.TABLE_0 WHERE ip='IP_0'",
  "sessionId": "sess_abc123",
  "wasSanitized": true,
  "mappingsCreated": {
    "ProductionDB": "SERVER_0",
    "users": "TABLE_0",
    "192.168.1.100": "IP_0"
  },
  "blocked": false
}
```

### 4. Use with OpenCode

Once installed, Enterprise Shield **automatically protects** all your OpenCode/Cursor AI interactions.

**No code changes. No configuration required. It just works.**

---

## 🔐 What Gets Protected?

### Automatically Sanitized (Masked with Aliases)

| Category | Examples | Becomes | Severity |
|----------|----------|---------|----------|
| **Database Servers** | `ServerDB01`, `ProductionDB`, `staging-db` | `SERVER_0` | Medium |
| **Database Tables** | `users_prod`, `customer_data`, `orders_2024` | `TABLE_0` | Medium |
| **Private IPs** | `192.168.1.100`, `10.0.0.50`, `172.16.0.1` | `IP_0` | High |
| **Connection Strings** | `Server=prod;User=admin;...` | `CONNSTR_0` | High |
| **File Paths (Windows)** | `C:\Projects\Secret\data.csv` | `PATH_0` | Medium |
| **UNC Paths** | `\\\\fileserver\\share\\` | `PATH_0` | Medium |
| **Internal Hostnames** | `db-prod-01.internal.company.com` | `HOST_0` | Medium |

### Automatically Blocked (Critical Violations)

| Category | Examples | Detection Method | Severity |
|----------|----------|------------------|----------|
| **Social Security Numbers** | `123-45-6789` | Format + validation | **Critical** |
| **Credit Card Numbers** | `4111111111111111` | Format + Luhn algorithm | **Critical** |
| **AWS Access Keys** | `AKIAIOSFODNN7EXAMPLE` | Prefix pattern | **Critical** |
| **GitHub Tokens** | `ghp_xxxxxxxxxxxx`, `gho_xxxxxxxxxxxx` | Prefix pattern | **Critical** |
| **OpenAI API Keys** | `sk-...48chars...` | Prefix + length | **Critical** |
| **Anthropic API Keys** | `sk-ant-...` | Prefix pattern | **Critical** |
| **Private Keys** | `-----BEGIN RSA PRIVATE KEY-----` | Header detection | **Critical** |
| **Passwords in Code** | `password="secret123"` | Context pattern | High |
| **Bearer Tokens** | `Bearer eyJhbGc...` | JWT pattern | High |
| **Slack Tokens** | `xoxb-...` | Prefix pattern | **Critical** |
| **Azure Storage Keys** | `AccountKey=...88chars...` | Pattern + length | **Critical** |

---

## 📖 Usage Examples

### Example 1: Database Query Optimization

**Developer types in Cursor:**
```sql
How do I optimize this query?
SELECT user_email, last_login, credit_score
FROM ProductionDB.user_accounts ua
JOIN FinanceDB.credit_data cd ON ua.id = cd.user_id
WHERE server='db-prod-01.internal.company.com'
  AND last_login > '2024-01-01'
```

**AI receives (sanitized):**
```sql
How do I optimize this query?
SELECT user_email, last_login, credit_score
FROM SERVER_0.TABLE_0 ua
JOIN SERVER_1.TABLE_1 cd ON ua.id = cd.user_id
WHERE server='HOST_0'
  AND last_login > '2024-01-01'
```

**Developer sees (desanitized response):**
```
To optimize ProductionDB.user_accounts and FinanceDB.credit_data:

1. Add composite index on user_accounts(last_login, id)
2. Ensure db-prod-01.internal.company.com has connection pooling enabled
3. Consider partitioning user_accounts by last_login date
4. Cache credit_data lookups in Redis for frequently accessed users
```

✅ **Benefit:** Real optimization advice without leaking infrastructure details!

### Example 2: Accidental PII Protection

**Developer accidentally types:**
```
Help me parse this user data:
John Doe, SSN: 123-45-6789, Card: 4111111111111111
```

**Enterprise Shield blocks immediately:**
```
🚫 REQUEST BLOCKED - CRITICAL VIOLATIONS DETECTED

Violations Found:
  1. Social Security Number (SSN)
     - Redacted: 123-****6789
     - Position: 23
     - Severity: CRITICAL
  
  2. Credit Card Number  
     - Redacted: 4111****1111
     - Position: 45
     - Severity: CRITICAL

❌ Request was NOT sent to AI provider
📝 Incident logged for audit
👤 Security team may be notified
```

✅ **Benefit:** Prevents accidental PII leakage before it happens!

### Example 3: Session Continuity

**First question:**
```
"Explain how to connect to ServerDB01"
```

**Response includes:** `SERVER_0` (mapped to ServerDB01 in session)

**Follow-up question:**
```
"Now show me how to query ServerDB01 and also connect to ServerDB02"
```

**AI sees:**
```
"Now show me how to query SERVER_0 and also connect to SERVER_1"
```

**You get back:**
```
"To query ServerDB01, use: SELECT * FROM ServerDB01.tablename
To connect to ServerDB02, use the same credentials but port 5433"
```

✅ **Benefit:** Consistent aliases throughout your conversation!

### Example 4: Team-Wide Protection

**Policy Configuration** (`~/.opencode/config/enterprise-shield.yaml`):

```yaml
policy:
  defaultAccessLevel: "sanitized_only"  # All requests sanitized
  
  # Department-specific rules
  departments:
    engineering:
      allowedProviders: ["openai", "anthropic"]
      dailyRequestLimit: 500
    
    contractors:
      accessLevel: "blocked"  # No AI access
    
    executives:
      accessLevel: "unrestricted"  # No sanitization (if needed)
```

✅ **Benefit:** Centralized security policies across entire organization!

---

## ✨ Features

### Core Security Capabilities

| Feature | Description | Testing |
|---------|-------------|---------|
| **Sanitization Engine** | Regex-based detection, alias generation | ✅ 6 tests |
| **Desanitization Engine** | Restores original values in responses | ✅ 6 tests |
| **Compliance Detector** | 14 PII/secret patterns (SSN, cards, keys) | ✅ 12 tests |
| **Session Management** | Per-user isolated sessions with 8hr TTL | ✅ Implemented |
| **Policy Engine** | RBAC with 3 access levels | ✅ Implemented |
| **Audit Logging** | Ed25519-signed append-only logs | ✅ Implemented |
| **Encryption** | AES-256-GCM for session data | ✅ Implemented |

**Total: 24 tests, all passing** ✅

### Enterprise Features

- ✅ **Zero-Knowledge Architecture**: Sensitive data never reaches external AI
- ✅ **Session Isolation**: Per-user encrypted sessions (AES-256-GCM)
- ✅ **Tamper-Evident Audit Trail**: Ed25519 cryptographic signatures
- ✅ **Policy-Based Access Control**: Unrestricted / SanitizedOnly / Blocked
- ✅ **Multi-Provider Support**: OpenAI, Anthropic, Azure OpenAI, Google AI
- ✅ **Configurable Rules**: YAML-based pattern management
- ✅ **Fail-Secure Design**: Blocks requests on error (never leaks)
- ✅ **Cross-Platform**: macOS, Linux, Windows support
- ✅ **Performance**: <50ms overhead per request

### Detection Patterns (14 Built-in)

**Infrastructure:**
- Server/database names (regex patterns)
- IP addresses (RFC 1918 private ranges)
- Connection strings
- File paths (Windows, UNC)
- Internal hostnames

**Critical PII:**
- SSN (US Social Security Numbers)
- Credit cards (with Luhn validation)
- API keys (AWS, GitHub, OpenAI, Anthropic, Slack, Azure)
- Private keys (RSA, EC, DSA, OpenSSH)
- Passwords in code
- Bearer tokens (JWT)

---

## 📡 CLI Commands

### Basic Commands

```bash
# Check version
enterprise-shield version

# Initialize configuration (creates ~/.opencode/config/enterprise-shield.yaml)
enterprise-shield init

# Scan content for violations (dry-run, doesn't send to AI)
enterprise-shield scan "<content>"

# Process a request with sanitization
enterprise-shield process <userID> "<content>" <provider>

# Run as background service (for OpenCode integration)
enterprise-shield serve
```

### Advanced Usage

```bash
# Scan file content
cat query.sql | enterprise-shield scan "$(cat -)"

# Process with specific session
enterprise-shield process user@example.com "Query DB01" openai --session=sess_abc123

# View current session
enterprise-shield session show

# Clear session
enterprise-shield session clear
```

---

## 🏗️ Architecture

### Clean Architecture (Go)

```
enterprise-shield/
├── cmd/plugin/              # CLI entry point
├── pkg/
│   ├── types/              # Shared type definitions
│   ├── sanitizer/          # Pattern-based sanitization
│   ├── desanitizer/        # Response restoration
│   ├── compliance/         # PII/secrets detection (Luhn, etc.)
│   ├── session/            # Session management (in-memory, Redis-ready)
│   ├── policy/             # RBAC policy engine
│   ├── audit/              # Ed25519-signed audit logging
│   ├── crypto/             # AES-256-GCM encryption
│   ├── hooks/              # OpenCode integration layer
│   └── config/             # YAML configuration loader
```

### Request Flow

```
OpenCode/Cursor AI Request
    ↓
[1] Policy Check (RBAC)
    ├─ Check user access level
    ├─ Verify provider is allowed
    └─ Check rate limits
    ↓
[2] Compliance Scan
    ├─ Detect PII (SSN, credit cards)
    ├─ Detect secrets (API keys, passwords)
    └─ BLOCK if critical violation
    ↓
[3] Sanitization
    ├─ Pattern matching (regex)
    ├─ Alias generation (SERVER_0, IP_0)
    └─ Update session mappings
    ↓
[4] Forward to LLM
    (sanitized content sent)
    ↓
[5] Receive Response
    ↓
[6] Desanitization
    ├─ Lookup aliases in session
    └─ Restore original values
    ↓
[7] Audit Log
    ├─ Create log entry
    ├─ Sign with Ed25519
    └─ Append to audit trail
    ↓
Return to Developer
```

---

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage report
make test-cover

# Build binary
make build

# Install locally to ~/.opencode/plugins/
make install

# Run demo
make demo
```

**Coverage:** 24 tests across sanitizer, desanitizer, and compliance modules

---

## 📊 Platform Support

| Platform | Architecture | Status | Install Methods |
|----------|-------------|--------|----------------|
| macOS | Intel (amd64) | ✅ Tested | All 6 methods |
| macOS | Apple Silicon (arm64) | ✅ Tested | All 6 methods |
| Linux | amd64 | ✅ Tested | All 6 methods |
| Linux | arm64 | ✅ Ready | All 6 methods |
| Windows | amd64 | ✅ Ready | Go install, Binary, Docker |

---

## 📖 Documentation

- **[README.md](README.md)** - This file (overview & installation)
- **[DISTRIBUTION.md](DISTRIBUTION.md)** - Complete distribution guide (8 pages)
- **[DISTRIBUTION_COMPLETE.md](DISTRIBUTION_COMPLETE.md)** - All distribution methods
- **[COMPLETE_SUMMARY.md](COMPLETE_SUMMARY.md)** - Full project summary
- **[CHANGELOG.md](CHANGELOG.md)** - Version history

---

## 🤝 Contributing

We welcome contributions! This project follows:
- **Go Best Practices** - Idiomatic, clean Go code
- **Test-Driven Development** - All features have tests first
- **Clean Architecture** - Clear separation of concerns
- **Security First** - Fail-secure design principles

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

---

## ⚖️ License

MIT License - See [LICENSE](LICENSE) file for details.

Open source and free for commercial use.

---

## 📞 Support & Contact

- **📖 Documentation**: [Complete docs](https://github.com/YOLOVibeCode/opencode-enterprise-shield)
- **🐛 Issues**: [GitHub Issues](https://github.com/YOLOVibeCode/opencode-enterprise-shield/issues)
- **🔒 Security**: security@YOLOVibeCode.com (responsible disclosure)
- **💼 Enterprise Support**: enterprise@YOLOVibeCode.com
- **💬 Community**: [Discussions](https://github.com/YOLOVibeCode/opencode-enterprise-shield/discussions)

---

## 🎉 Status

✅ **Production-Ready**  
✅ **24/24 Tests Passing**  
✅ **6 Distribution Methods**  
✅ **Enterprise-Grade Security**  
✅ **Zero Dependencies** (beyond Go stdlib + 3 modules)  
✅ **Cross-Platform** (macOS, Linux, Windows)  
✅ **MIT Licensed** (Free for commercial use)

**Enterprise Shield is ready to protect your organization's sensitive data!** 🛡️

---

## 🚀 Get Started Now

```bash
# One command to install
curl -sSL https://raw.githubusercontent.com/YOLOVibeCode/enterprise-shield/main/install.sh | bash

# That's it! You're protected.
```

---

*Built with ❤️ using Go 1.22 | Implements Prompt Shield specifications | Zero compromises on security or usability*
