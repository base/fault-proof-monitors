package tests

import (
	"fmt"
	"testing"
)

var (
	monitorEighteenFile = "incorrect_bond_balance.gate"
)

func TestIncorrectBondBalanceIncorrectFutureETHUnlocked(t *testing.T) {
	// We expect an alert to be fired when the FutureETHUnlocked is not the expected value

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
		"resolveClaimCalls": [][]interface{}{
			{3, 512},
			{2, 512},
			// claim indices 1 and 0 have not been resolved yet
		},
		"unlocksWithSender": [][]interface{}{
			// two unlocks for two resolveClaim calls
			{"0x00000000000000000000000000000000000000AA", []interface{}{"0x0000000000000000000000000000000000000001", 400}},
			{"0x00000000000000000000000000000000000000AA", []interface{}{"0x0000000000000000000000000000000000000002", 300}},
		},
		"ethBondsPerClaimIndex": []int{200, 100},   // the eth bonds for the remaining claim indices are 200 and 100
		"ethBondAtMinClaim":     0,                 // the min claim index, which is 2, is fully resolved
		"currDisputeEthBalance": 800,               // too much eth is remains to be unlocked, which will cause the alert to fire
		"pastWithdrawals":       [][]interface{}{}, // no past withdrawals have occurred
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

func TestIncorrectBondBalancePartiallyResolvedMinClaim(t *testing.T) {
	// We expect an alert to be fired when the FutureETHUnlocked is not the expected value due to a partially resolved min claim

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
		"resolveClaimCalls": [][]interface{}{
			{3, 512},
			{2, 1}, // claim 2 is only partially resolved, which means its bond has not been unlocked yet
			// claim indices 1 and 0 have not been resolved yet
		},
		"unlocksWithSender": [][]interface{}{
			// one unlock for one finished resolveClaim call
			{"0x00000000000000000000000000000000000000AA", []interface{}{"0x0000000000000000000000000000000000000001", 400}},
		},
		"ethBondsPerClaimIndex": []int{200, 100},   // the eth bonds for the remaining claim indices are 200 and 100
		"ethBondAtMinClaim":     300,               // the min claim index, which is 2, has not been fully resolved yet
		"currDisputeEthBalance": 800,               // too much eth is remains to be unlocked, which will cause the alert to fire
		"pastWithdrawals":       [][]interface{}{}, // no past withdrawals have occurred
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

func TestIncorrectBondBalanceIncorrectCurrETHUnlocked(t *testing.T) {
	// We expect an alert to be fired when the CurrentETHUnlocked is not the expected value

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
		"resolveClaimCalls": [][]interface{}{
			// all claims have been resolved
			{3, 512},
			{2, 512},
			{1, 512},
			{0, 512},
		},
		"unlocksWithSender": [][]interface{}{
			{"0x00000000000000000000000000000000000000AA", []interface{}{"0x0000000000000000000000000000000000000001", 400}},
			{"0x00000000000000000000000000000000000000AA", []interface{}{"0x0000000000000000000000000000000000000002", 300}},
			{"0x00000000000000000000000000000000000000AA", []interface{}{"0x0000000000000000000000000000000000000003", 200}},
			{"0x00000000000000000000000000000000000000AA", []interface{}{"0x0000000000000000000000000000000000000004", 100}},
		},
		"ethBondsPerClaimIndex": []int{},
		"ethBondAtMinClaim":     0,                 // all claims are fully resolved
		"currDisputeEthBalance": 800,               // too much eth has been unlocked, which will cause the alert to fire
		"pastWithdrawals":       [][]interface{}{}, // no past withdrawals have occurred
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
	// We DO NOT expect an alert to be fired when the FutureETHUnlocked and CurrentETHUnlocked are the expected values

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
		"resolveClaimCalls": [][]interface{}{
			{3, 512},
			{2, 512},
			// claim indices 1 and 0 have not been resolved yet
		},
		"unlocksWithSender": [][]interface{}{
			// two unlocks for two resolveClaim calls
			{"0x00000000000000000000000000000000000000AA", []interface{}{"0x0000000000000000000000000000000000000001", 400}},
			{"0x00000000000000000000000000000000000000AA", []interface{}{"0x0000000000000000000000000000000000000002", 300}},
			// this unlock is for a different dispute game so it should not be counted
			{"0x00000000000000000000000000000000000000BB", []interface{}{"0x0000000000000000000000000000000000000002", 300}},
		},
		"ethBondAtMinClaim":     0,               // the min claim index, which is 2, is fully resolved
		"ethBondsPerClaimIndex": []int{200, 100}, // the eth bonds for the remaining claim indices are 200 and 100, respectively
		"currDisputeEthBalance": 800,
		"currBlockNumber":       1,
		"pastWithdrawalEvents": [][]interface{}{
			// two withdrawal events have already happened, which when added to the current dispute
			// eth balance gives the correct amount
			{0, []interface{}{100}},
			{0, []interface{}{100}},
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

func TestIncorrectBondBalanceCorrectPastAndCurrentETHValues(t *testing.T) {
	// We DO NOT expect an alert to be fired when the FutureETHUnlocked and CurrentETHUnlocked are the expected values
	// where the CurrentETHUnlocked is aggregated from past and current withdrawals

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
		"resolveClaimCalls": [][]interface{}{
			{3, 512},
			{2, 512},
			// claim indices 1 and 0 have not been resolved yet
		},
		"unlocksWithSender": [][]interface{}{
			// two unlocks for two resolveClaim calls
			{"0x00000000000000000000000000000000000000AA", []interface{}{"0x0000000000000000000000000000000000000001", 400}},
			{"0x00000000000000000000000000000000000000AA", []interface{}{"0x0000000000000000000000000000000000000002", 300}},
			// this unlock is for a different dispute game so it should not be counted
			{"0x00000000000000000000000000000000000000BB", []interface{}{"0x0000000000000000000000000000000000000002", 300}},
		},
		"ethBondAtMinClaim":     0,               // the min claim index, which is 2, is fully resolved
		"ethBondsPerClaimIndex": []int{200, 100}, // the eth bonds for the remaining claim indices are 200 and 100, respectively
		"currDisputeEthBalance": 800,
		"currBlockNumber":       1,
		"pastWithdrawalEvents": [][]interface{}{
			// one withdrawal event has already happened, the other one happened in the latest block but will be accounted
			// for in the currentWithdrawalEvents (so no double count)
			{0, []interface{}{100}},
			{1, []interface{}{100}},
		},
		"currentWithdrawalEvents": [][]interface{}{
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

func TestIncorrectBondBalanceNoClaimsResolvedYet(t *testing.T) {
	// We DO NOT expect an alert to be fired when no claims have been resolved yet

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
		"claimIndices":     []int{}, // no claims have been resolved yet
		"unlockAmounts":    []int{},
		// when NO claims have been resolved, it is difficult to determine the depth of the remainin claim indices,
		// as the dispute game could still be going on - so instead when we haven't seen any claim yet, we just set
		// the value of the future eth unlocked to the currDisputeEthBalance
		"ethBondAtMinClaim":     0,       // no claims have been resolved yet
		"ethBondsPerClaimIndex": []int{}, // no indices are ready to be claimed yet
		"currDisputeEthBalance": 800,
		"pastWithdrawals":       [][]interface{}{}, // no past withdrawals have occurred
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
		"delayedWETH": "0x00000000000000000000000000000000000000BB",
		"resolveClaimCalls": [][]interface{}{
			{3, 512},
			{2, 512},
			// claim indices 1 and 0 have not been resolved yet
		},
		"unlocksWithSender": [][]interface{}{
			// two unlocks for two resolveClaim calls
			{"0x00000000000000000000000000000000000000AA", []interface{}{"0x0000000000000000000000000000000000000001", 400}},
			{"0x00000000000000000000000000000000000000AA", []interface{}{"0x0000000000000000000000000000000000000002", 300}},
		},
		"ethBondsPerClaimIndex": []int{200, 100},   // the eth bonds for the remaining claim indices are 200 and 100
		"ethBondAtMinClaim":     0,                 // the min claim index, which is 2, is fully resolved
		"currDisputeEthBalance": 800,               // too much eth is remains to be unlocked, which will cause the alert to fire
		"pastWithdrawals":       [][]interface{}{}, // no past withdrawals have occurred
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
