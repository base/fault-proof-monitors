package tests

import (
	"fmt"
	"testing"
)

var (
	monitorSeventeenFile = "credit_and_bond_discrepancy.gate"
)

func TestCreditAndBondDiscrepancyWrongBondAmount(t *testing.T) {
	// We expect an alert to be fired when the bond amount does not match the credit amount

	// set the param, which doesn't matter for this test suite
	params := map[string]any{
		"disputeGame": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorSeventeenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorSeventeenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"addressesInTrace": []any{"0x000000000000000000000000000000000000000"},
		"delayedWeth":      "0x0000000000000000000000000000000000000000", // doesn't matter
		"creditCalls": [][]any{
			// two addresses are claiming credit
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9"},
			{"0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4"},
		},
		"unlocks": [][]any{
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000000},
			{"0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", 1000000},
		},
		"withdraws":    [][]any{}, // no withdraw calls have been made
		"withdrawList": []any{},   // no withdraw calls have been made
		"winnersAndBonds": [][]any{
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000000},
			{"0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", 900000}, // does not match credit amount in unlocks
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorSeventeenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorSeventeenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorSeventeenFile)
	}
}

func TestCreditAndBondDiscrepancyMissingUnlock(t *testing.T) {
	// We expect an alert to be fired when there is a claimCredit call with no matchingunlock

	// set the param
	params := map[string]any{
		"disputeGame": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorSeventeenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorSeventeenFile, err)
	}

	// set the mock data
	mocks := map[string]any{
		"addressesInTrace": []any{"0x000000000000000000000000000000000000000"},
		"delayedWeth":      "0x0000000000000000000000000000000000000000",
		"creditCalls": [][]any{
			// two addresses are claiming credit
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9"},
			{"0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4"},
		},
		"unlocks": [][]any{
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000000}, // no unlock for the second address
		},
		"withdraws":    [][]any{}, // no withdraw calls have been made
		"withdrawList": []any{},   // no withdraw calls have been made
		"winnersAndBonds": [][]any{
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000000},
			{"0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", 1000000},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorSeventeenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorSeventeenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorSeventeenFile)
	}
}

func TestCreditAndBondDiscrepancyMissingWithdraw(t *testing.T) {
	// We expect an alert to be fired when there is a claimCredit call with no matching withdraw

	// set the param
	params := map[string]any{
		"disputeGame": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorSeventeenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorSeventeenFile, err)
	}

	// set the mock data
	mocks := map[string]any{
		"addressesInTrace": []any{"0x000000000000000000000000000000000000000"},
		"delayedWeth":      "0x0000000000000000000000000000000000000000",
		"creditCalls": [][]any{
			// two addresses are claiming credit
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9"},
			{"0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4"},
		},
		"unlocks": [][]any{}, // no unlocks are present
		"withdraws": [][]any{
			// one withdrawal is missing
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000000},
		},
		"winnersAndBonds": [][]any{},
		"foundUnlocks":    []bool{},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorSeventeenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorSeventeenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorSeventeenFile)
	}
}

func TestCreditAndBondDiscrepancyMissingWithdrawAndUnlock(t *testing.T) {
	// We expect an alert to be fired when there is a claimCredit call with no matching withdraw and unlock

	// set the param
	params := map[string]any{
		"disputeGame": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorSeventeenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorSeventeenFile, err)
	}

	// set the mock data
	mocks := map[string]any{
		"addressesInTrace": []any{"0x000000000000000000000000000000000000000"},
		"delayedWeth":      "0x0000000000000000000000000000000000000000",
		"creditCalls": [][]any{
			// two addresses are claiming credit
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9"},
			{"0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4"},
		},
		"unlocks":         [][]any{}, // no unlock calls have been made
		"withdraws":       [][]any{}, // no withdraw calls have been made
		"withdrawList":    []any{},
		"winnersAndBonds": [][]any{},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorSeventeenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorSeventeenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorSeventeenFile)
	}
}

func TestCreditAndBondDiscrepancyCorrectAmountsAndAddresses(t *testing.T) {
	// We DO NOT expect an alert to be fired if the bond and credit amounts match
	// and the claimant address matches the credited address

	// set the param
	params := map[string]any{
		"disputeGame": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorSeventeenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorSeventeenFile, err)
	}

	// set the mock data
	mocks := map[string]any{
		"addressesInTrace": []any{"0x000000000000000000000000000000000000000"},
		"delayedWeth":      "0x0000000000000000000000000000000000000000",
		"creditCalls": [][]any{
			// two addresses are claiming credit
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9"},
			{"0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4"},
		},
		"unlocks": [][]any{
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000000},
			{"0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", 1000000},
		},
		"withdraws": [][]any{
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000000},
			{"0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", 1000000},
		},
		"winnersAndBonds": [][]any{
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000000},
			{"0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", 1000000},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorSeventeenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorSeventeenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorSeventeenFile)
	}
}

func TestCreditAndBondDiscrepancyNoFilterAddress(t *testing.T) {
	// We DO NOT expect an alert to be fired when there is no address filtered in the current block trace

	// set the param, which doesn't matter for this test suite
	params := map[string]any{
		"disputeGame": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorSeventeenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorSeventeenFile, err)
	}

	// set the mock data
	mocks := map[string]any{
		"delayedWeth": "0x0000000000000000000000000000000000000000",
		"creditCalls": [][]any{
			// two addresses are claiming credit
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9"},
			{"0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4"},
		},
		"unlocks": [][]any{
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000000},
		},
		"withdraws":    [][]any{},
		"withdrawList": []any{},
		"winnersAndBonds": [][]any{
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000000},
			{"0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", 1000000},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorSeventeenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorSeventeenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorSeventeenFile)
	}
}
