# p/moul/authz

Composable authorization for Gno.

## TODO

- [x] transferable from static to dao-managed
- [x] easy to plug
- [x] easy to extend
- [x] english-first usage
- [ ] RBAC support (in addition to PBAC)
- [ ] shared interface between drivers

## Usage

```go
// One-line setup (call from init())
auth, priv := authz.NewWithOrigin()
safeConfig := priv.Safe() // Safe to expose

// Configure
priv.Add("propose", "g1user1", "g1user2")
priv.Add("vote", tokenHolders)

// Check
auth.AssertCurrentCan("propose")
auth.AssertPreviousCan("vote")

// Transfer control
safeConfig.TransferAdminTo(newMembership)
```

## Composable Memberships

```go
members := authz.And(
    authz.Or(users, tokenHolders),
    authz.Not(blacklist)
)
```

## Interfaces

- **Membership**: `Has(std.Address) bool`
- **Authority**: `Can()`, `AssertCan()`, `DoByCurrent()`, etc.
- **SafeConfig**: Exposable admin configuration

