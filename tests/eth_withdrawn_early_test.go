package tests

import (
	"fmt"
	"testing"
)

var (
	monitorTenFile = "eth_withdrawn_early.gate"
)

func TestETHWithdrawnTooEarly(t *testing.T) {
	// We expect an alert to be fired when a withdrawal is made before the delayedTime has passed

	// set the params
	params := map[string]any{
		"disputeGame": "0x00000000000000000000000000000000000000AA",
		"multicall3":  "0x00000000000000000000000000000000000000BB",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorTenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorTenFile, err)
	}

	// setup the mocks
	mocks := map[string]any{
		"addressesInTrace": []any{"0x00000000000000000000000000000000000000AA"},
		"delayedWETH":      "0x0000000000000000000000000000000000000000",
		"claims": [][]any{
			{"0x0000000000000000000000000000000000000001"},
		},
		// note claims and withdrawals are effectively the same thing, we just need to match claim calls
		// to their corresponding withdrawal call on a separate contract
		"withdrawals": [][]any{
			{"0x0000000000000000000000000000000000000001", 100},
		},
		"delayTime":        100,
		"currTimestamp":    1099,
		"unlockTimestamps": []int{1000, 1000}, // both unlocks happened earlier than the delayTime
		"unlocks": [][]any{
			// this unlock on its own would be fine as the delayTime has elapsed
			{50, "0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000001", 50}},
			// however, this unlock call is made before the delayTime has passed, and the delay resets each time a new unlock
			// is called against the same address for the same dispute game, so this will trigger the alert
			{101, "0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000001", 50}},
		},
		"hasUnlockedCredit": []any{true, true},
	}

	// call out to hexagate API to run the gate file with params and mocks
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorTenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorTenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorTenFile)
	}
}

func TestETHWithdrawnTooEarlyNoMatchingUnlock(t *testing.T) {
	// We expect an alert to be fired when a withdrawal is made but there is no matching unlock call
	// for the recipient address

	// set the params
	params := map[string]any{
		"disputeGame": "0x00000000000000000000000000000000000000AA",
		"multicall3":  "0x00000000000000000000000000000000000000BB",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorTenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorTenFile, err)
	}

	// setup the mocks
	mocks := map[string]any{
		"addressesInTrace": []any{"0x00000000000000000000000000000000000000AA"},
		"delayedWETH":      "0x0000000000000000000000000000000000000000",
		"claims": [][]any{
			{"0x0000000000000000000000000000000000000001"},
			{"0x0000000000000000000000000000000000000002"},
		},
		"withdrawals": [][]any{
			{"0x0000000000000000000000000000000000000001", 100},
			{"0x0000000000000000000000000000000000000002", 200},
		},
		"delayTime":        10,
		"currTimestamp":    2000,
		"unlockTimestamps": []int{1000}, // unlock happened after delayTime
		"unlocks": [][]any{
			{90, "0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000001", 100}},
			// we are missing the corresponding unlock for claim 0x00...02 which will trigger the alert
		},
		"hasUnlockedCredit": []any{true},
	}

	// call out to hexagate API to run the gate file with params and mocks
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorTenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorTenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorTenFile)
	}
}

func TestETHWithdrawnTooEarlyIncorrectAmount(t *testing.T) {
	// We expect an alert to be fired when a withdrawal is made but does not match the sum of the
	// unlock calls for the recipient address

	// set the params
	params := map[string]any{
		"disputeGame": "0x00000000000000000000000000000000000000AA",
		"multicall3":  "0x00000000000000000000000000000000000000BB",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorTenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorTenFile, err)
	}

	// setup the mocks
	mocks := map[string]any{
		"addressesInTrace": []any{"0x00000000000000000000000000000000000000AA"},
		"delayedWETH":      "0x0000000000000000000000000000000000000000",
		"claims": [][]any{
			{"0x0000000000000000000000000000000000000001"},
			{"0x0000000000000000000000000000000000000002"},
		},
		"withdrawals": [][]any{
			{"0x0000000000000000000000000000000000000001", 100},
			{"0x0000000000000000000000000000000000000002", 200},
		},
		"delayTime":        10,
		"currTimestamp":    2000,
		"unlockTimestamps": []int{1000, 1000}, // unlocks happened after delayTime
		"unlocks": [][]any{
			{90, "0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000001", 100}},
			// notice the unlock amount for claim 0x00...02 doesn't match the withdrawal amount
			{90, "0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000002", 100}},
		},
		"hasUnlockedCredit": []any{true, true}, // credit was unlocked for both claims
	}

	// call out to hexagate API to run the gate file with params and mocks
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorTenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorTenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorTenFile)
	}
}

func TestETHWithdrawnTooEarlyNoUnlockedCredit(t *testing.T) {
	// We expect an alert to be fired when a withdrawal is made but the recipient has not unlocked their credit

	// set the params
	params := map[string]any{
		"disputeGame": "0x00000000000000000000000000000000000000AA",
		"multicall3":  "0x00000000000000000000000000000000000000BB",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorTenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorTenFile, err)
	}

	// setup the mocks
	mocks := map[string]any{
		"addressesInTrace": []any{"0x00000000000000000000000000000000000000AA"},
		"delayedWETH":      "0x0000000000000000000000000000000000000000",
		"claims": [][]any{
			{"0x0000000000000000000000000000000000000001"},
			{"0x0000000000000000000000000000000000000002"},
		},
		"withdrawals": [][]any{
			{"0x0000000000000000000000000000000000000001", 100},
			{"0x0000000000000000000000000000000000000002", 200},
		},
		"delayTime":        10,
		"currTimestamp":    2000,
		"unlockTimestamps": []int{1000, 1000, 1000}, // unlocks happened after delayTime
		"unlocks": [][]any{
			{90, "0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000001", 100}},
			// aggregating the two unlocks together will give us the correct withdrawal amount
			{90, "0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000002", 100}},
			{89, "0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000002", 100}},
		},
		"hasUnlockedCredit": []any{true, false}, // only one claim has unlocked credit
	}

	// call out to hexagate API to run the gate file with params and mocks
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorTenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorTenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorTenFile)
	}
}

func TestETHWithdrawnTooEarlyCorrectWithdrawal(t *testing.T) {
	// We DO NOT expect an alert to be fired when a withdrawal occurrs past the delayedTime,
	// with the correct sum and matching unlock calls

	// set the params
	params := map[string]any{
		"disputeGame": "0x00000000000000000000000000000000000000AA",
		"multicall3":  "0x00000000000000000000000000000000000000BB",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorTenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorTenFile, err)
	}

	// setup the mocks
	mocks := map[string]any{
		"addressesInTrace": []any{"0x00000000000000000000000000000000000000AA"},
		"delayedWETH":      "0x0000000000000000000000000000000000000000",
		"claims": [][]any{
			{"0x0000000000000000000000000000000000000001"},
			{"0x0000000000000000000000000000000000000002"},
		},
		"withdrawals": [][]any{
			{"0x0000000000000000000000000000000000000001", 100},
			{"0x0000000000000000000000000000000000000002", 200},
		},
		"delayTime":        10,
		"currTimestamp":    2000,
		"unlockTimestamps": []int{1000, 1000, 1000}, // unlocks happened after delayTime
		"unlocks": [][]any{
			{90, "0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000001", 100}},
			// aggregating the two unlocks together will give us the correct withdrawal amount
			{90, "0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000002", 100}},
			{89, "0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000002", 100}},
		},
		"hasUnlockedCredit": []bool{true, true},
	}

	// call out to hexagate API to run the gate file with params and mocks
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorTenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorTenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorTenFile)
	}
}

func TestETHWithdrawnTooEarlyNoClaimInBlock(t *testing.T) {
	// We DO NOT expect an alert to be fired when there is no claim in the current block
	// set the params
	params := map[string]any{
		"disputeGame": "0x00000000000000000000000000000000000000AA",
		"multicall3":  "0x00000000000000000000000000000000000000CC",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorTenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorTenFile, err)
	}

	// setup the mocks
	mocks := map[string]any{
		"addressesInTrace": []any{"0x00000000000000000000000000000000000000AA"},
		"delayedWETH":      "0x0000000000000000000000000000000000000000",
		"claims":           [][]any{},
		"withdrawals": [][]any{
			// does not have matching claim in this block so this will get parsed out
			{"0x0000000000000000000000000000000000000001", 200},
		},
		"delayTime":        10,
		"currTimestamp":    2000,
		"unlockTimestamps": []int{1000}, // unlock happened after delayTime
		"unlocks": [][]any{
			{90, "0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000001", 100}},
			{90, "0x00000000000000000000000000000000000000BB", []any{"0x0000000000000000000000000000000000000002", 100}},
		},
		"hasUnlockedCredit": []bool{true},
	}

	// call out to hexagate API to run the gate file with params and mocks
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorTenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorTenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorTenFile)
	}
}

func TestETHWithdrawnTooEarlyNoFilterAddress(t *testing.T) {
	// We DO NOT expect an alert to be fired when there is no address in the filter trace

	// set the params
	params := map[string]any{
		"disputeGame": "0x00000000000000000000000000000000000000AA",
		"multicall3":  "0x00000000000000000000000000000000000000BB",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorTenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorTenFile, err)
	}

	// setup the mocks
	mocks := map[string]any{
		"delayedWETH": "0x0000000000000000000000000000000000000000",
		"claims": [][]any{
			{"0x0000000000000000000000000000000000000001"},
		},
		// note claims and withdrawals are effectively the same thing, we just need to match claim calls
		// to their corresponding withdrawal call on a separate contract
		"withdrawals": [][]any{
			{"0x0000000000000000000000000000000000000001", 100},
		},
		"delayTime":        100,
		"currTimestamp":    1099,
		"unlockTimestamps": []int{1000, 1000}, // both unlocks happened earlier than the delayTime
		"unlocks": [][]any{
			// this unlock on its own would be fine as the delayTime has elapsed
			{50, "0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000001", 50}},
			// however, this unlock call is made before the delayTime has passed, and the delay resets each time a new unlock
			// is called against the same address for the same dispute game, so this will trigger the alert
			{101, "0x00000000000000000000000000000000000000AA", []any{"0x0000000000000000000000000000000000000001", 50}},
		},
		"hasUnlockedCredit": []bool{true},
	}

	// call out to hexagate API to run the gate file with params and mocks
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorTenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorTenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorTenFile)
	}
}
