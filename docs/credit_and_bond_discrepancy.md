## Purpose

The `credit_and_bond_discrepancy.gate` monitor ensures the **safety** and integrity of the fault proof system by verifying that the ETH bond amounts credited during the resolution of dispute games match the unlock amounts. The monitor also verifies that any call that would interact with the two-step bond withdrawal process occurs with either an unlock or withdraw call. Discrepancies indicate potential bugs in the dispute game contracts or the contracts that hold the bonds (`DelayedWETH`), which could compromise the incentivization mechanism and put user funds at risk.

## Technical Overview

### How It Works

1. **Monitoring Credit Claims**: The monitor listens for `claimCredit` function calls on the dispute game contract, which indicate that a participant is requesting their credit.

2. **Retrieving Bond Unlock Calls**:

   - The monitor retrieves all `unlock` function calls on the `DelayedWETH` contract (the contract holding the bonds), which indicate bonds being unlocked and credited to participants.

3. **Retrieving Withdraw Calls**:

   - The monitor retrieves all `withdraw` function calls on the `DelayedWETH` contract (the contract holding the bonds), which indicate bonds being withdrawn and sent to participants.

4. **Determining Credit Amount Per Recipient**:

     - The view function `credit` is called for all `claimCredit` recipients to determine the expected unlock amount.

5. **Cross-Referencing Values**:

   - The monitor cross-references the expected recipients and credit amounts with the actual unlocks recorded.
   - It checks that for each expected recipient and amount, there is a matching unlock event.
   - The monitor also cross-references the expected recipients and withdrawals, as these should also occur together.

6. **Detecting Discrepancies**:

   - If any expected unlock is not found in the actual `unlock` calls, it indicates a discrepancy.
   - If any expected withdraw is not found in the actual `withdraw` calls, it indicates a discrepancy.
   - If neither an unlock or withdraw occurs when a `claimCredit` call occurs, this also indicates a discrepancy.

7. **Triggering Alerts**:

   - If discrepancies are detected, the monitor raises an alert for immediate investigation.

### Importance of the Monitor

- **Ensuring Correct Incentivization**: The dispute game mechanism relies on proper financial incentives. Participants are motivated to act honestly because they risk losing their bonds if they behave maliciously. Over or under-accounting of ETH bonds undermines this mechanism.

- **Preventing Financial Loss**: Discrepancies in bond amounts can lead to financial losses for participants who do not receive the correct bond value upon resolution.

- **Maintaining Trust in the System**: Verifying the two-step bond withdrawal process is essential for participants to trust the dispute resolution process and continue participating, which is crucial for the fault proof system's functionality.

- **Detecting Critical Bugs**: Discrepancies may indicate bugs or logic flaws in the dispute game contracts or related components, requiring prompt attention to prevent further issues.

## Parameters

- `disputeGame`: Address of the dispute game contract being monitored.