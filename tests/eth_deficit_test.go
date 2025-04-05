package tests

import (
	"fmt"
	"testing"
)

var (
	monitorElevenFile = "eth_deficit.gate"
)

func TestETHDeficitTotalCreditDeficitNormalMode(t *testing.T) {
	// We expect an alert to be fired when totalCredit is less than claimCredit and bondDistributionMode is NORMAL

	// set the params, which don't really matter for these tests
	params := map[string]any{
		"disputeGame":      "0x0000000000000000000000000000000000000000",
		"honestChallenger": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorElevenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorElevenFile, err)
	}

	// set the mock data
	mocks := map[string]any{
		"delayedWETH":           "0x0000000000000000000000000000000000000000",
		"bondDistributionMode":  1,    // NORMAL mode, so claimCredit will be used in the invariant check
		"hasUnlockedCredit":     true, // challenger has unlocked credit
		"claimCredit":           100,
		"refundModeCredit":      0,                 // doesn't matter in this test
		"totalCredit":           []any{50, 123456}, // (amount, timestamp)
		"ethBalanceDisputeGame": 500,
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorElevenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorElevenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorElevenFile)
	}
}

func TestETHDeficitTotalCreditDeficitRefundMode(t *testing.T) {
	// We expect an alert to be fired when totalCredit is less than refundModeCredit and bondDistributionMode is REFUND

	// set the params, which don't really matter for these tests
	params := map[string]any{
		"disputeGame":      "0x0000000000000000000000000000000000000000",
		"honestChallenger": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorElevenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorElevenFile, err)
	}

	// set the mock data
	mocks := map[string]any{
		"delayedWETH":           "0x0000000000000000000000000000000000000000",
		"bondDistributionMode":  2,    // REFUND mode, so refundModeCredit will be used in the invariant check
		"hasUnlockedCredit":     true, // challenger has unlocked credit
		"claimCredit":           0,    // doesn't matter in this test
		"refundModeCredit":      100,
		"totalCredit":           []any{50, 123456}, // (amount, timestamp)
		"ethBalanceDisputeGame": 500,
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorElevenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorElevenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorElevenFile)
	}
}

func TestETHDeficitTotalETHBalanceDeficit(t *testing.T) {
	// We expect an alert to be fired when ethBalanceDisputeGame is less than totalCredit

	// set the params
	params := map[string]any{
		"disputeGame":      "0x0000000000000000000000000000000000000000",
		"honestChallenger": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorElevenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorElevenFile, err)
	}

	// set the mock data
	mocks := map[string]any{
		"delayedWETH":           "0x0000000000000000000000000000000000000000",
		"bondDistributionMode":  1,    // just needs to be 1 or 2
		"hasUnlockedCredit":     true, // challenger has unlocked credit
		"claimCredit":           50,   // less than totalCredit so the invariant won't fire on this check
		"refundModeCredit":      50,
		"totalCredit":           []any{150, 123456},
		"ethBalanceDisputeGame": 100,
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorElevenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorElevenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorElevenFile)
	}
}

func TestETHDeficitCurrCreditZeroTotalCreditNonZero(t *testing.T) {
	// We expect an alert to be fired when claimCredit is zero and totalCredit is non-zero

	// set the params
	params := map[string]any{
		"disputeGame":      "0x0000000000000000000000000000000000000000",
		"honestChallenger": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorElevenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorElevenFile, err)
	}

	// set the mock data
	mocks := map[string]any{
		"delayedWETH":           "0x0000000000000000000000000000000000000000",
		"bondDistributionMode":  1,    // just needs to be 1 or 2
		"hasUnlockedCredit":     true, // challenger has unlocked credit
		"claimCredit":           0,
		"refundModeCredit":      0,
		"totalCredit":           []any{150, 123456},
		"ethBalanceDisputeGame": 1500,
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorElevenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorElevenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorElevenFile)
	}
}

func TestETHDeficitNoCreditUnlocked(t *testing.T) {
	// We expect an alert to be fired when the amount of credit unlocked is non-zero but unlockedCredit is false

	// set the params, which don't really matter for these tests
	params := map[string]any{
		"disputeGame":      "0x0000000000000000000000000000000000000000",
		"honestChallenger": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorElevenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorElevenFile, err)
	}

	// set the mock data
	mocks := map[string]any{
		"delayedWETH":           "0x0000000000000000000000000000000000000000",
		"bondDistributionMode":  1,     // just needs to be 1 or 2
		"hasUnlockedCredit":     false, // challenger has NOT unlocked credit
		"claimCredit":           50,    // value needs to be less than or equal to totalCredit
		"refundModeCredit":      50,
		"totalCredit":           []any{50, 123456},
		"ethBalanceDisputeGame": 500,
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorElevenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorElevenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorElevenFile)
	}
}

func TestETHDeficitBondDistributionModeZero(t *testing.T) {
	// We DO NOT expect an alert to be fired if bondDistributionMode is 0 (game undecided)

	// set the params, which don't really matter for these tests
	params := map[string]any{
		"disputeGame":      "0x0000000000000000000000000000000000000000",
		"honestChallenger": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorElevenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorElevenFile, err)
	}

	// set the mock data
	mocks := map[string]any{
		"delayedWETH":           "0x0000000000000000000000000000000000000000",
		"bondDistributionMode":  0,     // undecided game, so the invariant should not fire
		"hasUnlockedCredit":     false, // credit has not been unlocked yet
		"claimCredit":           0,     // has not been decided yet
		"refundModeCredit":      50,    // increments whenever a move occurs, so will be nonzero
		"totalCredit":           []any{50, 123456},
		"ethBalanceDisputeGame": 500,
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorElevenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorElevenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorElevenFile)
	}
}

func TestETHDeficitNoDeficit(t *testing.T) {
	// We DO NOT expect an alert to be fired if there is no deficit

	// set the params
	params := map[string]any{
		"disputeGame":      "0x0000000000000000000000000000000000000000",
		"honestChallenger": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorElevenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorElevenFile, err)
	}

	// set the mock data
	mocks := map[string]any{
		"delayedWETH":           "0x0000000000000000000000000000000000000000",
		"bondDistributionMode":  1,    // just needs to be 1 or 2
		"hasUnlockedCredit":     true, // challenger has unlocked credit
		"claimCredit":           50,   // value needs to be less than or equal to totalCredit
		"refundModeCredit":      50,
		"totalCredit":           []any{50, 123456},
		"ethBalanceDisputeGame": 500,
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorElevenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorElevenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorElevenFile)
	}
}
