## Purpose

The `incorrect_bond_balance.gate` monitor ensures the **safety** and integrity of the fault proof system by verifying that the total ETH bonds associated with a dispute game are correctly accounted for during the resolution process. Specifically, it checks that the sum of ETH bonds that have been unlocked is at most equal to the maximum amount of ETH that the dispute game has ever held. Additionally, the maximum amount of ETH held by the dispute game is also checked for equivalence against the total bonds aggregated from all claims. Any imbalance indicates potential issues such as multiple resolutions or over-accounted claims, which can compromise the incentivization mechanism and result in financial losses for participants.

## Technical Overview

### How It Works

1. **Monitoring Credit Unlocks**:

   - Retrieves all historical `unlock` calls on the `DelayedWETH` contract that originated from the current `disputeGame`.

2. **Calculating Total Bond Value of all Claims in the Dispute Game**:

   - Retrieves the bond values of all the claims in the `disputeGame` and sums them.

3. **Calculating Total ETH Held by the Dispute Game**:

   - Retrieves the current ETH balance of the `disputeGame` in the `DelayedWETH` contract.
   - Sums the current ETH balance with any past ETH withdrawals to determine the maximum amount of ETH that the dispute game has ever held.

4. **Validating Bond Accounting**:

   - Verifies that the sum of the bond amounts equals the total ETH that has ever been held by the dispute game.
   - Verifies that the sum of credit unlocks is less than or equal to the total ETH that has ever been held by the dispute game.
   - Checks for any discrepancies, which would indicate an imbalance.

5. **Triggering Alerts**:

   - If an imbalance is detected, the monitor raises an alert for immediate investigation.

### Importance of the Monitor

- **Ensuring Correct Incentivization**: Accurate bond accounting is essential for the incentivization mechanism of dispute games. Under-accounting or over-accounting of bonds undermines trust and the proper functioning of the system.

- **Preventing Financial Loss**: Imbalances can lead to financial losses for participants who are entitled to receive bond amounts upon dispute resolution.

- **Detecting Critical Issues**: Discrepancies may indicate that subgames have been skipped, resolved out of order, or resolved multiple times, pointing to broader issues in the dispute game logic.

## Parameters

- `disputeGame`: Address of the dispute game contract being monitored.
