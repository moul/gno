package dao_sudo

import (
	"std"
	"testing"
	"time"
)

func TestDAOContractDeployment(t *testing.T) {
	// Setup: Add members to DAO
	owner := std.GetOrigCaller()
	member1 := std.Address("g1member1xxxxxxxxxxxxxxxxxxxxxxxxxxx")
	member2 := std.Address("g1member2xxxxxxxxxxxxxxxxxxxxxxxxxxx")
	
	// Add members with voting power
	AddMember(member1, 30)
	AddMember(member2, 30)
	// Owner already has 100, so total is 160
	
	// Create a proposal to deploy a new contract
	contractCode := `package hello

import "std"

var message = "Hello from DAO-deployed contract!"

func GetMessage() string {
    return message
}

func SetMessage(newMsg string) {
    caller := std.GetOrigCaller()
    if caller != std.Address("` + std.GetRealmAddr().String() + `") {
        panic("only DAO can update message")
    }
    message = newMsg
}`
	
	proposalData := "gno.land/r/demo/dao_deployed/hello|" + contractCode
	
	// Create deployment proposal
	proposalID := CreateProposal(
		"Deploy Hello Contract",
		"Deploy a simple hello world contract controlled by the DAO",
		ProposalTypeDeployContract,
		proposalData,
	)
	
	// Members vote on the proposal
	std.TestSetOrigCaller(member1)
	Vote(proposalID, true) // 30 yes votes
	
	std.TestSetOrigCaller(member2)
	Vote(proposalID, true) // 30 yes votes
	
	std.TestSetOrigCaller(owner)
	Vote(proposalID, true) // 100 yes votes
	
	// Total: 160 yes votes out of 160 total (100% support)
	
	// Fast forward past voting period
	// In real test, we'd use a test helper to advance time
	proposal := getProposal(proposalID)
	proposal.EndTime = time.Now().Add(-1 * time.Hour)
	
	// Execute the proposal
	ExecuteProposal(proposalID)
	
	// Verify proposal was executed
	if proposal.Status != ProposalStatusExecuted {
		t.Errorf("Expected proposal status to be Executed, got %v", proposal.Status)
	}
	
	// The contract deployment would be queued via realm.Sudo
	// and executed after this transaction completes
}

func TestDAOContractCall(t *testing.T) {
	owner := std.GetOrigCaller()
	
	// Create a proposal to call a contract function
	proposalData := "gno.land/r/demo/dao_deployed/hello|SetMessage|Hello from DAO!"
	
	proposalID := CreateProposal(
		"Update Hello Message",
		"Update the message in our deployed contract",
		ProposalTypeCallContract,
		proposalData,
	)
	
	// Vote to pass the proposal
	Vote(proposalID, true)
	
	// Fast forward and execute
	proposal := getProposal(proposalID)
	proposal.EndTime = time.Now().Add(-1 * time.Hour)
	
	ExecuteProposal(proposalID)
	
	// Verify execution
	if proposal.Status != ProposalStatusExecuted {
		t.Errorf("Expected proposal status to be Executed, got %v", proposal.Status)
	}
}