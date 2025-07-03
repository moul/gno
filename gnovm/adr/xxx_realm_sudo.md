# ADR-XXX: Contract Transactions (realm.Sudo)

## Status

Proposed

## Context

Smart contracts need to generate new transactions that execute after their own completion. Key use cases:
- DAOs deploying approved contracts
- Contracts sending tokens or executing bank operations
- Multi-step workflows requiring separate transaction contexts
- Contracts triggering any type of blockchain action

## Decision

Implement `realm.Sudo(msg Message)` to queue messages for execution after the current transaction completes. Messages are routed through the standard message router with the contract as sender.

```go
package dao

func ExecuteProposal(proposal Proposal) {
    switch proposal.Type {
    case "deploy":
        m := vm.NewMsgAddPkg(GetRealmAddr(), proposal.PkgPath, proposal.Code)
        realm.Sudo(m)
    case "send":
        m := bank.NewMsgSend(GetRealmAddr(), proposal.Recipient, proposal.Amount)
        realm.Sudo(m)
    case "delegate":
        m := staking.NewMsgDelegate(GetRealmAddr(), proposal.Validator, proposal.Amount)
        realm.Sudo(m)
    }
}
```

## Implementation

1. **realm.Sudo** - Queues any message type (only callable by realms)
2. **Message Router** - Routes messages to appropriate handlers after main tx
3. **Context** - Contract address automatically set as message sender
4. **Security** - Recursion limits, gas accounting per message

## Consequences

**Pros**: Enables DAOs, maintains tx boundaries, preserves determinism
**Cons**: Increased complexity, requires careful gas management

## References

- Cosmos SDK message handling
- Example: `/examples/gno.land/r/demo/dao_sudo/`