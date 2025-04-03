## Purpose

The `eth_withdrawn_early.gate` monitor ensures the **safety** and integrity of the fault proof system by detecting any instances where ETH withdrawals from the `DelayedWETH` contract occur earlier than the defined delay period (`DELAY_SECONDS`). The delay period is crucial as it allows privileged actors time to intervene if an error or unexpected behavior is detected in a dispute game. Bypassing this delay compromises the dispute resolution process and can lead to unauthorized or premature withdrawals, potentially risking significant amounts of ETH bonds held in the `DelayedWETH` contract.

## Technical Overview

### How It Works

1. **Monitoring Withdrawals**:

   - **Claim Credit Calls**: The monitor observes calls to the `claimCredit` function on the `disputeGame` contract, which are initiated by participants to claim their ETH credits after a dispute resolution.

   - **Withdraw Calls**: It tracks corresponding `withdraw` function calls on the `DelayedWETH` contract to capture when ETH is actually withdrawn by participants.

2. **Correlating Unlocks and Withdrawals**:

   - **Unlock Events**: Retrieves historical `unlock` function calls on the `DelayedWETH` contract that are initiated by the specified `disputeGame`. These events record when ETH bonds are unlocked and become eligible for withdrawal after the delay period.

   - **Timestamps**: Uses the `Multicall3` contract to obtain historical block timestamps for both unlock and withdraw events to calculate the time difference between them.

3. **Validating Delay Period**:

   - **Delay Enforcement**: Checks that the time elapsed between the `unlock` event and the corresponding `withdraw` event is at least equal to the `DELAY_SECONDS` specified in the `DelayedWETH` contract.

4. **Ensuring Amount Consistency**:

   - **Amount Verification**: Verifies that the amount withdrawn matches the total amount unlocked for the recipient, ensuring that participants cannot withdraw more than they are entitled to.

5. **Ensuring Credit Unlock Status**:

   - **Credit Unlock Status**: Verifies that for each withdrawal recipient, their credit unlock status is `true` in the specified `disputeGame`.

6. **Triggering Alerts**:

   - **Early Withdrawal Detection**: If a withdrawal occurs before the required delay period has elapsed, or if the withdrawn amount does not match the unlocked amount, the monitor raises an alert for immediate investigation.

### Importance of the Monitor

- **Preventing Unauthorized Withdrawals**: Ensures that participants cannot bypass the delay mechanism and withdraw ETH prematurely, which could indicate malicious activity or exploitation of a vulnerability.

- **Maintaining System Integrity**: The enforced delay period is a critical safeguard, allowing time for privileged actors to respond to any issues in dispute games. Bypassing this period compromises the security protocols and trust in the network.

- **Protecting User Funds**: Unauthorized or early withdrawals can lead to significant financial losses, especially since the `DelayedWETH` contract holds the cumulative bonds from all active dispute games.

## Parameters

- `multicall3`: Address of the `Multicall3` contract used to retrieve block timestamps.
- `disputeGame`: Address of the dispute game contract being monitored.
