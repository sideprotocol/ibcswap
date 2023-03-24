---
title: Deposit(lock)
stage: draft
category: IBCSWAP/Interchainswap
kind: instantiation
author: Marian <marian@side.one>
created: 2023-3-24
modified: 2023-3-24
---

## Synopsis

Upon activating a pool, it is crucial to safeguard against spam deposits and implement a deposit timing logic flow to shield the pool from potential collapse. A malicious actor may attempt to deposit a large sum and withdraw it rapidly to manipulate the asset price within a short time frame.

### Motivation

This algorithm is designed to defend the pool against spam deposit/withdraw operations, thus contributing to the stabilization of asset prices within the pool.

### Definitions

- `Lock Ticks`: This term refers to the timeline points for the lock duration of deposited tokens within the pool, indicating the number of days the deposit will remain locked. Currently planed three options `1,14,28 days`. If user deposit, he can't withdraw token before expiring lock date.

## Technical Specification

### Algorithm

#### Plan 1

##### Scenario

In this approach, we will create a deposit ledger for each depositor, storing the LP token amount and lock expiration date based on the deposit request.

When a user attempts to withdraw their funds, we will identify all eligible withdrawal deposit requests and calculate the available amount. If the available amount is greater than the requested withdrawal, we will process the withdrawal request and update the sender's ledger accordingly.

- define deposit metadata `DepositMetadata`:

      ```
        enum DepositLockTicks {
          LOCK7  = 0;
          LOCK14 = 1;
          LOCK28 = 2;
        }

        type DepositMetadata struct {
            ID uint64
            PoolToken sdk.Coin
            IssuedDate time.Time
            Tick DepositLockTicks
        }
      ```

- create `depositLedger`. `depositLeger` is a combined data type `map` and `list`
  ```
    deopositLedger := make(map[sdk.Address][]DepositMetadata{})
  ```
- update `OnDepositReceived`

  ```
    if pool.Status == types.PoolStatus_POOL_STATUS_READY {
  	// update pool tokens.
  	for _, token := range msg.Tokens {
  		pool.AddAsset(*token)
  	}
  } else {
  	// switch pool status to 'READY'
  	pool.Status = types.PoolStatus_POOL_STATUS_READY
  }

    newDepositID := len(depositLedger[msg.sender])
    newDeposit := DepositMetadata {
        ID: newDepositID,
        PoolToken: poolToken,
        ExpireDate: ctx.BlockTime(),
    }

    k.AddNewDepositToLedger(ctx, msg.sender, newDeposit)
  ```

  - update `OnSingleDepositAcknowledged`

  ```
  func (k Keeper) OnSingleDepositAcknowledged(ctx sdk.Context, req *types.MsgDepositRequest, res *types.MsgDepositResponse) error {
    ...
    newDepositID := len(depositLedger[req.sender])
    newDeposit := DepositMetadata {
        ID: newDepositID,
        PoolToken: *res.PoolToken,
        ExpireDate ctx.BlockTime() + msg.Tick,
    }
    k.AddNewDepositToLedger(ctx, msg.sender, newDeposit)
  }
  ```

  - update withdraw

  ```
   func (k Keeper) OndWithdrawReceive(ctx sdk.Context, msg *types.MsgWithdrawRequest) (*types.MsgWithdrawResponse, error) {

  	if err := msg.ValidateBasic(); err != nil {
  		return nil, err
  	}
  	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolCoin.Denom)

  	if !found {
  		return nil, types.ErrNotFoundPool
  	}

      depositLedger := GetDepositLedger(ctx, msg.sender)
      withdrawableAmount := 0
      startIndex := 0
      for index, metadata := range depositLedger {
          if metadata.ExpireDate > time.now() {
              withdrawableAmount.add(metadata.PoolToken)
          }else{
              break
          }
          if withdrawableAmount >=  msg.Tokens[0] {
              startIndex = index
              remainAmount := withdrawableAmount.sub(metadata.PoolToken)
              if remainAmount.Uint64() > sdk.NewInt(0) {
                  depositLedger[index].Tokens = remainAmount
                   startIndex = index
              }else{
                  startIndex = index+1
              }
              break
          }
      }

    if withdrawableAmount < msg.Token[0] {
        return nil, ErrNotEnoughWithdrawableTokens
        }

      if startIndex != 0 {
          updatedDepositLedger = depositLedger[startIndex:len(depositLedger)]
          k.SetDepositLedger(ctx, msg.sender, updatedDepositLedger)
      }
      ...
  }
  ```

Add the same update logic `OnWithdrawAcknowledged`.

##### Note

- In this case, users obtain their LP tokens immediately within their wallets.
- However, there is a loop to calculate the withdrawable amount. If the ledger continues to grow, it may impact the chain's performance. Implementing a limitation on the number of deposits may be a better approach. We can assume that a user is a malicious actor if they do not withdraw at all and only attempt to deposit multiple times. Additionally, we should establish a minimum deposit amount to prevent attacks.

#### Plan 2

##### Scenario

In this method, we do not directly provide Pool Tokens to users. Instead, we supply an NFT that includes specific metadata.

When a user deposits funds, their LP tokens will be locked in an escrow account, and the module will mint an NFT containing DepositMetadata.

##### Advantage

- Users do not hold Pool Tokens directly, preventing malicious withdrawal attempts.
- If users wish to withdraw, they must claim LP tokens from escrow instead of the NFT, eliminating complex ledger management within the protocol and shifting responsibility to the user.
- Additional services can be built on top of the NFT (e.g., secondary and tertiary vesting, voting, etc.)

##### Drawback

- Users must first obtain their LP tokens from the locker, requiring an additional step for withdrawal.
- It appears to compromise generalization as a widely-used IBC module.
- Minting NFTs necessitates the inclusion of the cosmos/x/nft module, increasing dependencies.
- On the frontend, users must manually select available deposits for withdrawal from their deposit NFT list, which may slightly reduce user experience.
