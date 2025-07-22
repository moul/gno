package sessions

import (
	"strings"
	"time"
)

// This file provides comprehensive documentation and examples for the account sessions feature

// Example code snippets for using sessions
var examples = []struct {
	Title       string
	Description string
	Code        string
}{
	{
		Title:       "Create a Daily Hot Wallet Session",
		Description: "Create a session with limited daily spending capacity",
		Code: `gnokey session create \
  -account alice \
  -name daily-wallet \
  -transfer-capacity 1000000ugnot \
  -expire-in 24h \
  -gas-wanted 1000000 \
  -gas-fee 1ugnot \
  -broadcast`,
	},
	{
		Title:       "Create a DAO Voting Session",
		Description: "Create a session restricted to governance contracts",
		Code: `gnokey session create \
  -account alice \
  -name dao-voter \
  -realms "gno.land/r/demo/dao/*,gno.land/r/demo/governance" \
  -gas-wanted 1000000 \
  -gas-fee 1ugnot \
  -broadcast`,
	},
	{
		Title:       "Create a Validator Session",
		Description: "Create a validation-only session for node operations",
		Code: `gnokey session create \
  -account validator \
  -name validator-key \
  -validation-only \
  -gas-wanted 1000000 \
  -gas-fee 1ugnot \
  -broadcast`,
	},
	{
		Title:       "Create an Admin Session",
		Description: "Create a session that can manage other sessions",
		Code: `gnokey session create \
  -account alice \
  -name admin-session \
  -can-manage-sessions \
  -transfer-capacity 5000000ugnot \
  -gas-wanted 1000000 \
  -gas-fee 1ugnot \
  -broadcast`,
	},
	{
		Title:       "Use a Session for Transfers",
		Description: "Send funds using a session key",
		Code: `gnokey maketx send \
  -send 100000ugnot \
  -to g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5 \
  -gas-wanted 50000 \
  -gas-fee 1ugnot \
  -broadcast daily-wallet`,
	},
	{
		Title:       "Use a Session for Contract Calls",
		Description: "Call a realm function using a session key",
		Code: `gnokey maketx call \
  -pkgpath gno.land/r/demo/dao \
  -func Vote \
  -args "proposal-1" "yes" \
  -gas-wanted 100000 \
  -gas-fee 1ugnot \
  -broadcast dao-voter`,
	},
	{
		Title:       "List Account Sessions",
		Description: "View all sessions for an account",
		Code: `gnokey session list -account alice`,
	},
	{
		Title:       "Revoke a Session",
		Description: "Remove a session's permissions",
		Code: `gnokey session revoke \
  -account alice \
  -session-key-name daily-wallet \
  -gas-wanted 100000 \
  -gas-fee 1ugnot \
  -broadcast`,
	},
}

// RenderDocs returns documentation based on the path
func RenderDocs(path string) string {
	parts := strings.Split(path, "/")
	
	switch {
	case path == "" || path == "docs":
		return renderDocsHome()
	case path == "docs/examples":
		return renderExamples()
	case path == "docs/features":
		return renderFeatures()
	case path == "docs/security":
		return renderSecurityGuide()
	case path == "docs/patterns":
		return renderPatterns()
	default:
		return "Documentation not found for path: " + path
	}
}

func renderDocsHome() string {
	return `# Account Sessions Documentation

Account sessions provide a secure way to delegate limited permissions to different keys, enabling advanced key management patterns.

## Documentation Sections

- [Examples](/r/demo/sessions:docs/examples) - Code examples for common operations
- [Features](/r/demo/sessions:docs/features) - Detailed feature overview
- [Security](/r/demo/sessions:docs/security) - Security best practices
- [Patterns](/r/demo/sessions:docs/patterns) - Common usage patterns

## Why Use Sessions?

1. **Security**: Keep your master key offline while using session keys for daily operations
2. **Automation**: Create restricted keys for bots and automated systems
3. **Delegation**: Share limited access without exposing your main key
4. **Compliance**: Implement organizational key management policies

## Quick Start

1. Create a session with specific permissions
2. Use the session key for allowed operations
3. Monitor session usage
4. Revoke when no longer needed

See the [examples](/r/demo/sessions:docs/examples) section for detailed commands.
`
}

func renderExamples() string {
	var output strings.Builder
	output.WriteString("# Session Command Examples\n\n")
	
	for _, example := range examples {
		output.WriteString("## ")
		output.WriteString(example.Title)
		output.WriteString("\n\n")
		output.WriteString(example.Description)
		output.WriteString("\n\n```bash\n")
		output.WriteString(example.Code)
		output.WriteString("\n```\n\n")
	}
	
	return output.String()
}

func renderFeatures() string {
	return `# Session Features

## Core Capabilities

### 1. Transfer Capacity Limits
- Set maximum transfer amount per session
- Capacity decreases with each use
- Prevents unlimited spending if key is compromised
- Master key always has unlimited capacity

### 2. Realm Access Control
- Restrict which contracts a session can interact with
- Use glob patterns for flexible matching
- Examples:
  - ` + "`gno.land/r/demo/*`" + ` - All demo realms
  - ` + "`gno.land/r/demo/dao,gno.land/r/demo/gov`" + ` - Specific realms
  - Empty list allows access to all realms

### 3. Time-Based Expiration
- Set absolute expiration time
- Use relative durations (1h, 24h, 7d)
- Expired sessions automatically become invalid
- Garbage collection removes expired sessions

### 4. Permission Flags
- **Unlimited Transfer**: No capacity restrictions
- **Manage Sessions**: Create/revoke other sessions
- **Manage Packages**: Deploy new contracts
- **Validation Only**: Restrict to validator operations

### 5. Independent Sequence Numbers
- Each key has its own sequence counter
- Prevents replay attacks
- Enables parallel transactions
- Account tracks total sequences

## Technical Details

### Session Structure
- Name: Human-readable identifier
- Public Key: Unique key for the session
- Sequence: Transaction counter
- Expiration: Optional time limit
- Transfer Capacity: Spending limit
- Realm Whitelist: Access restrictions
- Flags: Permission settings

### Garbage Collection
- Automatic cleanup of expired sessions
- Happens during account operations
- Configurable per blockchain
- Maintains account efficiency

### Key Relationships
- One key can control multiple accounts
- One account can have multiple keys
- Maximum 64 sessions per account
- Master key retains full control
`
}

func renderSecurityGuide() string {
	return `# Security Best Practices

## Key Storage

### Master Key
- **Cold Storage**: Keep offline, ideally hardware wallet
- **Limited Use**: Only for critical operations
- **Backup**: Multiple secure backups
- **Never Expose**: Never on internet-connected systems

### Session Keys
- **Purpose-Specific**: One key per use case
- **Secure Generation**: Automatic cryptographic entropy
- **Appropriate Storage**: Match security to permissions
- **Regular Rotation**: Replace keys periodically

## Permission Design

### Principle of Least Privilege
1. Grant minimum necessary permissions
2. Use transfer limits appropriate to use case
3. Restrict realm access when possible
4. Set reasonable expiration times

### Risk-Based Approach
- **High Value**: Shorter expiration, lower limits
- **Automation**: Strict realm restrictions
- **Human Use**: Balance security and usability
- **Validators**: Validation-only flag

## Operational Security

### Monitoring
- Track transfer capacity usage
- Monitor realm interactions
- Alert on unusual patterns
- Regular permission audits

### Incident Response
1. **Detection**: Monitor for anomalies
2. **Revocation**: Immediate key removal
3. **Investigation**: Analyze transaction history
4. **Remediation**: Update security measures

### Key Rotation Policy
- **Regular Schedule**: Based on risk level
- **Event-Driven**: After personnel changes
- **Compromise**: Immediate rotation
- **Audit Trail**: Document all changes

## Common Vulnerabilities

### Avoid These Mistakes
1. **Overpermissioning**: Too many capabilities
2. **No Expiration**: Permanent sessions
3. **Shared Keys**: Multiple users per key
4. **Poor Storage**: Insecure key storage
5. **No Monitoring**: Blind to misuse

### Secure Patterns
1. **Time-Boxing**: Always set expiration
2. **Capacity Limits**: Reasonable maximums
3. **Realm Restrictions**: Minimize access
4. **Regular Audits**: Review active sessions
5. **Incident Plans**: Prepared responses
`
}

func renderPatterns() string {
	return `# Common Usage Patterns

## 1. Hot Wallet Pattern

**Use Case**: Daily operations with limited risk

\` + "``" + `
- Transfer Capacity: 1-5% of total assets
- Expiration: 24-48 hours
- Realms: Commonly used contracts
- Storage: Encrypted on device
\` + "``" + `

**Benefits**:
- Convenient for regular use
- Limited loss if compromised
- Easy to rotate frequently

## 2. Service Account Pattern

**Use Case**: Automated systems and bots

\` + "``" + `
- Transfer Capacity: Minimal for fees
- Expiration: Long-term or none
- Realms: Specific contracts only
- Flags: No session management
\` + "``" + `

**Benefits**:
- Automated operations
- Restricted permissions
- Audit trail of actions

## 3. Validator Pattern

**Use Case**: Node operations

\` + "``" + `
- Transfer Capacity: None
- Expiration: Long-term
- Flags: Validation only
- Storage: On validator node
\` + "``" + `

**Benefits**:
- Secure validation
- No transfer risk
- Simple permission model

## 4. Delegation Pattern

**Use Case**: Temporary access for others

\` + "``" + `
- Transfer Capacity: Task-specific
- Expiration: Short-term
- Realms: Required contracts
- Monitoring: Active alerts
\` + "``" + `

**Benefits**:
- No key sharing
- Time-limited access
- Full audit trail
- Easy revocation

## 5. Multi-Sig Alternative

**Use Case**: Organizational control

\` + "``" + `
- Multiple sessions with different permissions
- Each user has their own session
- Admin session for management
- Regular audit and rotation
\` + "``" + `

**Benefits**:
- Individual accountability
- Flexible permissions
- No coordination overhead
- Simpler than multi-sig

## Implementation Tips

1. **Start Small**: Begin with one pattern
2. **Document Everything**: Clear procedures
3. **Train Users**: Proper key handling
4. **Monitor Actively**: Set up alerts
5. **Review Regularly**: Adjust as needed

## Emergency Procedures

### Key Compromise
1. Revoke affected session immediately
2. Review transaction history
3. Rotate related keys
4. Update security measures

### Lost Session Key
1. Revoke from master account
2. Create new session
3. Update dependent systems
4. Document incident

### Expired Sessions
1. Plan rotation in advance
2. Create new before expiration
3. Update automation
4. Remove old keys from storage
`
}