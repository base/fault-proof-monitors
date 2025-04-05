## Purpose

The `eth_deficit.gate` monitor ensures the **safety** and integrity of the fault proof system by detecting any deficits of ETH in the `DelayedWETH` contract associated with the specified dispute game. The `DelayedWETH` contract holds the ETH bonds deposited by participants in dispute games. A deficit indicates that more ETH has been withdrawn than should be allowed, potentially due to bugs in bond accounting or dispute game resolution, which can result in financial losses for honest participants. Note this monitor tracks deficits in relation to the Challenger (`honestChallenger`), who is assumed to be operating honestly and participating in every dispute game as necessary.

## Technical Overview

### How It Works

1. **Retrieving Key Balances and Credits**:

   - **DelayedWETH Address**: Retrieves the `DelayedWETH` contract address associated with the specific `disputeGame`.

   - **Challenger's Normal Mode Credit (`claimCredit`)**: The amount of ETH that the Challenger (`honestChallenger`) can currently claim based on the resolution of the `disputeGame`.

   - **Challenger's Refund Mode Credit (`refundModeCredit`)**: The amount of ETH that the Challenger can claim based on the moves made during the `disputeGame`.

   - **Bond Distribution Mode (`bondDistributionMode`)**: Determines the credit amounts participants will receive from the `disputeGame`.

   - **Challenger's Unlocked Credit Status (`hasUnlockedCredit`)**: Determines whether the Challenger has successfully unlocked their credit, which is a requirement before it can be withdrawn.

   - **Challenger's Total Credit (`totalCredit`)**: The total amount of ETH that has been unlocked for the `honestChallenger` in the `DelayedWETH` contract.

   - **Dispute Game's ETH Balance (`ethBalanceDisputeGame`)**: The total amount of ETH held in the `DelayedWETH` contract for the `disputeGame`.

2. **Validating Balances**:

   - **Credit Unlock Status**:
     - Ensures that the challenger has unlocked their credit using the `hasUnlockedCredit` value.

   - **Credit Consistency**:
     - Ensures that the credit that the challenger can claim, based on the `bondDistributionMode`, does not exceed the `totalCredit` (what has been unlocked for them).
     - Verifies that the amount the challenger is trying to claim is consistent with what is available.

   - **Total Credit vs. Dispute Game Balance**:
     - Ensures that the `totalCredit` does not exceed the `ethBalanceDisputeGame`.
     - Checks that the total credits unlocked for participants do not exceed the ETH actually held for the dispute game.

   - **Synchronization Check**:
     - Ensures that if the credit the challenger can claim is zero, then `totalCredit` should also be zero.
     - Detects discrepancies that might indicate desynchronization between the dispute game and the `DelayedWETH` contract.

3. **Triggering Alerts**:

   - If any of the above conditions fail, the monitor raises an alert indicating a potential deficit or inconsistency in the ETH balances related to the dispute game.

### Importance of the Monitor

- **Preventing Financial Loss**: A deficit in the `DelayedWETH` contract can lead to losses for honest participants expecting to receive their bonds back upon dispute resolution.

- **Ensuring Correct Bond Accounting**: Accurate tracking of bonds is crucial for the incentivization mechanism of dispute games. Over or under-accounting undermines trust and the proper functioning of the system.

- **Detecting Critical Issues Early**: Prompt identification of discrepancies allows for immediate investigation and correction of potential bugs in bond accounting or resolution logic.

- **Maintaining System Integrity**: Ensures that the dispute game mechanism operates securely, preserving the safety and reliability of the network.

## Parameters

- `disputeGame`: Address of the dispute game contract being monitored.
- `honestChallenger`: Address of the Challenger.
