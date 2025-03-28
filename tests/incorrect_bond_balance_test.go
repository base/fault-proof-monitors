package tests

import (
	"fmt"
	"testing"
)

var (
	monitorEighteenFile = "incorrect_bond_balance.gate"
)

func TestIncorrectBondBalanceIncorrectTotalDisputeETH(t *testing.T) {
	// We expect an alert to be fired when the totalDisputeEthBalance value is not equal to the totalClaimBonds value

	// set the params
	params := map[string]any{
		"disputeGame": "0x00000000000000000000000000000000000000AA",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorEighteenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorEighteenFile, err)
	}

	// setup the mocks
	mocks := map[string]any{
		"addressesInTrace": []any{"0x00000000000000000000000000000000000000AA"},
		"delayedWETH":      "0x00000000000000000000000000000000000000BB",
		"unlockAmounts":    []any{0}, // no unlocks have occurred
		// the only value we take from claimData is at index 3, which is the bond value for the claim
		"claimData": [][]any{
			{1, "0x00000000000000000000000000000000000000AA", "0x00000000000000000000000000000000000000AA", 100, "0x00", 1, 1},
		},
		"currDisputeEthBalance": 200,       // the balance of ETH in the DelayedWETH contract for the dispute game exceeds the total claim bonds
		"pastWithdrawals":       [][]any{}, // no past withdrawals have occurred
	}

	// call out to hexagate API to run the gate file with params and mocks
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorEighteenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorEighteenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorEighteenFile)
	}
}

func TestIncorrectBondBalanceIncorrectETHUnlocked(t *testing.T) {
	// We expect an alert to be fired when the total bonds unlocked exceeds the totalDisputeEthBalance

	// set the params
	params := map[string]any{
		"disputeGame": "0x00000000000000000000000000000000000000AA",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorEighteenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorEighteenFile, err)
	}

	// setup the mocks
	mocks := map[string]any{
		"addressesInTrace": []any{"0x00000000000000000000000000000000000000AA"},
		"delayedWETH":      "0x00000000000000000000000000000000000000BB",
		"unlocksWithSender": [][]any{
			// multiple unlocks have happened that exceed the total eth balance for the dispute game
			{"0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000001", 200}},
			{"0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000002", 200}},
			{"0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000003", 200}},
		},
		// the only value we take from claimData is at index 3, which is the bond value for the claim
		"claimData": [][]any{
			{1, "0x00000000000000000000000000000000000000AA", "0x00000000000000000000000000000000000000AA", 100, "0x00", 1, 1},
			{2, "0x00000000000000000000000000000000000000AA", "0x00000000000000000000000000000000000000AA", 100, "0x00", 1, 1},
		},
		"currDisputeEthBalance": 200,       // the balance of ETH in the DelayedWETH contract for the dispute game is equal to the total claim bonds
		"pastWithdrawals":       [][]any{}, // no past withdrawals have occurred
	}

	// call out to hexagate API to run the gate file with params and mocks
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorEighteenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorEighteenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorEighteenFile)
	}
}

func TestIncorrectBondBalanceIncorrectWithdrawalValue(t *testing.T) {
	// We expect an alert to be fired when the totalClaimBonds value is not equal to the totalDisputeEthBalance when past withdrawals are taken into account

	// set the params
	params := map[string]any{
		"disputeGame": "0x00000000000000000000000000000000000000AA",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorEighteenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorEighteenFile, err)
	}

	// setup the mocks
	mocks := map[string]any{
		"addressesInTrace": []any{"0x00000000000000000000000000000000000000AA"},
		"delayedWETH":      "0x00000000000000000000000000000000000000BB",
		"unlocksWithSender": [][]any{
			// a single unlock has happened that is equal to the total claim bonds
			{"0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000001", 200}},
		},
		// the only value we take from claimData is at index 3, which is the bond value for the claim
		"claimData": [][]any{
			{1, "0x00000000000000000000000000000000000000AA", "0x00000000000000000000000000000000000000AA", 100, "0x00", 1, 1},
			{2, "0x00000000000000000000000000000000000000AA", "0x00000000000000000000000000000000000000AA", 100, "0x00", 1, 1},
		},
		"currDisputeEthBalance": 200, // the balance of ETH in the DelayedWETH contract for the dispute game is equal to the total claim bonds
		"pastWithdrawalEvents": [][]any{
			// two withdrawal events have already happened, which when added to the current dispute
			// eth balance do not add up to the total claim bonds
			{100},
			{100},
		},
	}

	// call out to hexagate API to run the gate file with params and mocks
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorEighteenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorEighteenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorEighteenFile)
	}
}

func TestIncorrectBondBalanceCorrectETHValues(t *testing.T) {
	// We DO NOT expect an alert to be fired when the totalDisputeEthBalance, totalClaimBonds, and total unlocks are all equal

	// set the params
	params := map[string]any{
		"disputeGame": "0x00000000000000000000000000000000000000AA",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorEighteenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorEighteenFile, err)
	}

	// setup the mocks
	mocks := map[string]any{
		"addressesInTrace": []any{"0x00000000000000000000000000000000000000AA"},
		"delayedWETH":      "0x00000000000000000000000000000000000000BB",
		"unlocksWithSender": [][]any{
			// four unlocks have happened that are equal to the total claim bonds
			{"0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000001", 100}},
			{"0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000002", 100}},
			{"0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000003", 100}},
			{"0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000004", 100}},
		},
		// the only value we take from claimData is at index 3, which is the bond value for the claim
		"claimData": [][]any{
			{1, "0x00000000000000000000000000000000000000AA", "0x00000000000000000000000000000000000000AA", 100, "0x00", 1, 1},
			{2, "0x00000000000000000000000000000000000000AA", "0x00000000000000000000000000000000000000AA", 100, "0x00", 1, 1},
			{3, "0x00000000000000000000000000000000000000AA", "0x00000000000000000000000000000000000000AA", 100, "0x00", 1, 1},
			{4, "0x00000000000000000000000000000000000000AA", "0x00000000000000000000000000000000000000AA", 100, "0x00", 1, 1},
		},
		"currDisputeEthBalance": 200, // the balance of ETH in the DelayedWETH contract for the dispute game is equal to the total claim bonds
		"pastWithdrawalEvents": [][]any{
			// two withdrawal events have already happened, which when added to the current dispute
			// eth balance do not add up to the total claim bonds
			{100},
			{100},
		},
	}

	// call out to hexagate API to run the gate file with params and mocks
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorEighteenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorEighteenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorEighteenFile)
	}
}

func TestIncorrectBondBalanceNoFilterAddress(t *testing.T) {
	// We DO NOT expect an alert to be fired when the filter address is not in the trace

	// set the params
	params := map[string]any{
		"disputeGame": "0x00000000000000000000000000000000000000AA",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorEighteenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorEighteenFile, err)
	}

	// setup the mocks
	mocks := map[string]any{
		"addressesInTrace": []any{}, // No applicable addresses found, so even if the other mocks would cause the invariant to break, the monitor should not fire
		"delayedWETH":      "0x00000000000000000000000000000000000000BB",
		"unlocksWithSender": [][]any{
			// a single unlock has happened that is equal to the total claim bonds
			{"0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000001", 200}},
		},
		// the only value we take from claimData is at index 3, which is the bond value for the claim
		"claimData": [][]any{
			{1, "0x00000000000000000000000000000000000000AA", "0x00000000000000000000000000000000000000AA", 100, "0x00", 1, 1},
			{2, "0x00000000000000000000000000000000000000AA", "0x00000000000000000000000000000000000000AA", 100, "0x00", 1, 1},
		},
		"currDisputeEthBalance": 200, // the balance of ETH in the DelayedWETH contract for the dispute game is equal to the total claim bonds
		"pastWithdrawalEvents": [][]any{
			// two withdrawal events have already happened, which when added to the current dispute
			// eth balance do not add up to the total claim bonds
			{100},
			{100},
		},
	}

	// call out to hexagate API to run the gate file with params and mocks
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorEighteenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorEighteenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorEighteenFile)
	}
}
