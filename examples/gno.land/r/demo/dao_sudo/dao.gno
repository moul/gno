// Package dao_sudo demonstrates a DAO that can deploy new contracts
// after governance approval using realm.Sudo functionality.
package dao_sudo

import (
	"std"
	"strings"
	"time"
	
	"gno.land/p/demo/avl"
	"gno.land/p/demo/ufmt"
	"realm"
)

// Proposal represents a governance proposal
type Proposal struct {
	ID          uint64
	Title       string
	Description string
	Type        ProposalType
	Data        string // Contract code for deployment proposals
	Creator     std.Address
	CreatedAt   time.Time
	EndTime     time.Time
	YesVotes    uint64
	NoVotes     uint64
	Status      ProposalStatus
	Voters      *avl.Tree // address -> bool (true if voted)
}

type ProposalType uint8

const (
	ProposalTypeText ProposalType = iota
	ProposalTypeDeployContract
	ProposalTypeCallContract
	ProposalTypeSendTokens
	ProposalTypeDelegate
)

type ProposalStatus uint8

const (
	ProposalStatusActive ProposalStatus = iota
	ProposalStatusPassed
	ProposalStatusFailed
	ProposalStatusExecuted
)

var (
	proposals  *avl.Tree // ID -> *Proposal
	members    *avl.Tree // address -> voting power (uint64)
	nextID     uint64
	quorum     uint64 = 51 // 51% quorum required
	minPeriod  = 24 * time.Hour
	owner      std.Address
)

func init() {
	proposals = avl.NewTree()
	members = avl.NewTree()
	owner = std.GetOrigCaller()
	
	// Add initial member
	members.Set(owner.String(), uint64(100))
}

// AddMember adds a new DAO member with voting power
func AddMember(addr std.Address, votingPower uint64) {
	assertOwner()
	members.Set(addr.String(), votingPower)
}

// CreateProposal creates a new governance proposal
func CreateProposal(title, description string, proposalType ProposalType, data string) uint64 {
	caller := std.GetOrigCaller()
	if !isMember(caller) {
		panic("only DAO members can create proposals")
	}
	
	id := nextID
	nextID++
	
	proposal := &Proposal{
		ID:          id,
		Title:       title,
		Description: description,
		Type:        proposalType,
		Data:        data,
		Creator:     caller,
		CreatedAt:   time.Now(),
		EndTime:     time.Now().Add(minPeriod),
		Voters:      avl.NewTree(),
		Status:      ProposalStatusActive,
	}
	
	proposals.Set(ufmt.Sprintf("%d", id), proposal)
	return id
}

// Vote casts a vote on a proposal
func Vote(proposalID uint64, support bool) {
	caller := std.GetOrigCaller()
	votingPower := getVotingPower(caller)
	if votingPower == 0 {
		panic("not a DAO member")
	}
	
	proposal := getProposal(proposalID)
	if proposal.Status != ProposalStatusActive {
		panic("proposal is not active")
	}
	
	if time.Now().After(proposal.EndTime) {
		finalizeProposal(proposal)
		panic("voting period has ended")
	}
	
	// Check if already voted
	if _, exists := proposal.Voters.Get(caller.String()); exists {
		panic("already voted on this proposal")
	}
	
	// Record vote
	proposal.Voters.Set(caller.String(), true)
	if support {
		proposal.YesVotes += votingPower
	} else {
		proposal.NoVotes += votingPower
	}
}

// ExecuteProposal executes a passed proposal
func ExecuteProposal(proposalID uint64) {
	proposal := getProposal(proposalID)
	
	// Finalize if voting period ended
	if time.Now().After(proposal.EndTime) && proposal.Status == ProposalStatusActive {
		finalizeProposal(proposal)
	}
	
	if proposal.Status != ProposalStatusPassed {
		panic("proposal has not passed")
	}
	
	switch proposal.Type {
	case ProposalTypeDeployContract:
		executeContractDeployment(proposal)
	case ProposalTypeCallContract:
		executeContractCall(proposal)
	case ProposalTypeSendTokens:
		executeSendTokens(proposal)
	case ProposalTypeDelegate:
		executeDelegate(proposal)
	default:
		// Text proposals don't need execution
		proposal.Status = ProposalStatusExecuted
	}
}

// executeContractDeployment uses realm.Sudo to deploy a new contract
func executeContractDeployment(proposal *Proposal) {
	// Parse deployment data
	// Expected format: "pkgpath|code"
	parts := strings.Split(proposal.Data, "|")
	if len(parts) != 2 {
		panic("invalid deployment data format")
	}
	
	pkgPath := parts[0]
	code := parts[1]
	
	// Create a MsgAddPackage to deploy the contract
	// The VM keeper will automatically set the DAO's address as Creator
	// when processing this sudo message
	msg := &MsgAddPackage{
		Creator: std.GetRealmAddr(), // DAO's address will be set by VM keeper
		Package: Package{
			Name: extractPackageName(pkgPath),
			Path: pkgPath,
			Files: []MemFile{
				{
					Name: "contract.gno",
					Body: code,
				},
			},
		},
		Deposit: std.Coins{}, // Can add deposit if needed
	}
	
	// Queue the message for execution after this transaction
	// The message will be processed with the DAO's address as sender
	realm.Sudo(msg)
	
	proposal.Status = ProposalStatusExecuted
}

// executeContractCall uses realm.Sudo to call another contract
func executeContractCall(proposal *Proposal) {
	// Parse call data
	// Expected format: "pkgpath|function|args"
	parts := strings.Split(proposal.Data, "|")
	if len(parts) != 3 {
		panic("invalid call data format")
	}
	
	pkgPath := parts[0]
	function := parts[1]
	args := strings.Split(parts[2], ",")
	
	// Create a MsgCall to call the contract
	// The context will have the DAO's address as sender
	msg := &MsgCall{
		Caller:  std.GetRealmAddr(), // DAO's address  
		PkgPath: pkgPath,
		Func:    function,
		Args:    args,
		Send:    std.Coins{}, // Can add coins if needed
	}
	
	// Queue the message for execution after this transaction
	realm.Sudo(msg)
	
	proposal.Status = ProposalStatusExecuted
}

// executeSendTokens uses realm.Sudo to send tokens from the DAO
func executeSendTokens(proposal *Proposal) {
	// Parse send data
	// Expected format: "recipient|amount|denom"
	parts := strings.Split(proposal.Data, "|")
	if len(parts) != 3 {
		panic("invalid send data format")
	}
	
	recipient := std.Address(parts[0])
	amount := std.MustParseCoins(parts[1] + parts[2])
	
	// Create a MsgSend to transfer tokens
	msg := &MsgSend{
		FromAddress: std.GetRealmAddr(),
		ToAddress:   recipient,
		Amount:      amount,
	}
	
	// Queue the message for execution
	realm.Sudo(msg)
	
	proposal.Status = ProposalStatusExecuted
}

// executeDelegate uses realm.Sudo to delegate tokens
func executeDelegate(proposal *Proposal) {
	// Parse delegation data
	// Expected format: "validator|amount|denom"
	parts := strings.Split(proposal.Data, "|")
	if len(parts) != 3 {
		panic("invalid delegation data format")
	}
	
	validator := std.Address(parts[0])
	amount := std.MustParseCoins(parts[1] + parts[2])
	
	// Create a MsgDelegate
	msg := &MsgDelegate{
		DelegatorAddress: std.GetRealmAddr(),
		ValidatorAddress: validator,
		Amount:           amount,
	}
	
	// Queue the message for execution
	realm.Sudo(msg)
	
	proposal.Status = ProposalStatusExecuted
}

// Helper functions

func finalizeProposal(proposal *Proposal) {
	totalPower := getTotalVotingPower()
	quorumReached := (proposal.YesVotes + proposal.NoVotes) >= (totalPower * quorum / 100)
	
	if quorumReached && proposal.YesVotes > proposal.NoVotes {
		proposal.Status = ProposalStatusPassed
	} else {
		proposal.Status = ProposalStatusFailed
	}
}

func getProposal(id uint64) *Proposal {
	key := ufmt.Sprintf("%d", id)
	value, exists := proposals.Get(key)
	if !exists {
		panic("proposal not found")
	}
	return value.(*Proposal)
}

func isMember(addr std.Address) bool {
	_, exists := members.Get(addr.String())
	return exists
}

func getVotingPower(addr std.Address) uint64 {
	value, exists := members.Get(addr.String())
	if !exists {
		return 0
	}
	return value.(uint64)
}

func getTotalVotingPower() uint64 {
	total := uint64(0)
	members.Iterate("", "", func(key string, value interface{}) bool {
		total += value.(uint64)
		return false
	})
	return total
}

func extractPackageName(pkgPath string) string {
	parts := strings.Split(pkgPath, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "contract"
}

func assertOwner() {
	if std.GetOrigCaller() != owner {
		panic("unauthorized")
	}
}

// Render function for displaying DAO state
func Render(path string) string {
	if path == "" {
		return renderHome()
	}
	
	if strings.HasPrefix(path, "proposal/") {
		idStr := strings.TrimPrefix(path, "proposal/")
		// Parse ID and render proposal
		return renderProposal(idStr)
	}
	
	return "404: Not Found"
}

func renderHome() string {
	output := "# DAO with Contract Deployment\n\n"
	output += "This DAO can deploy new contracts after governance approval.\n\n"
	
	output += "## Active Proposals\n\n"
	
	proposals.Iterate("", "", func(key string, value interface{}) bool {
		proposal := value.(*Proposal)
		if proposal.Status == ProposalStatusActive {
			output += ufmt.Sprintf("- [%d: %s](/r/demo/dao_sudo:proposal/%d)\n", 
				proposal.ID, proposal.Title, proposal.ID)
		}
		return false
	})
	
	return output
}

func renderProposal(idStr string) string {
	// Implementation details omitted for brevity
	return "Proposal details..."
}