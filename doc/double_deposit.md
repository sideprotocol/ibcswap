Similar to ICA, Sending an ibc packet from chain A, include signature of chainB user's signature and execute Tx on chain B according to the signature verification result.

a tx structure will similar to this:

```js
{
     senders: ["Achain1xxxxxxxxx", "Bchain1xxxxx"],
     tokens: [AchainCoin, BchainCoin],
     signature: []bytes
}
```

steps:

1. construct a tx on chain A
   - build the Counter party chain sender's signature.
   ```js
   depositTx := &types.EncounterPartyDepositTx{
   				AccountSequence: 1,
   				Sender:          chainBAddress,
   				Token:           &sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(initialY)},
   	}
    rawDepositTx := types.CDC.MustMarshal(depositTx)
    signedTx := priv.Sign(rawDepositTx)
   ```
   - build double deposit message.
   ```js
   msg := types.NewMsgDoubleDeposit(
   				poolId: string,
   				senders: string[],
   				depositTokens: sdk.Coins[],
   				signedTx: bytes,
   )
   ```
2. sending to from chain A
3. relayer tx
4. received tx on Chain B, build message(banktypes.MsgSend) based on packet information and verify signature from Bchain sender.

```js
    acc := k.accountKeeper.GetAccount(ctx, msg.senders[1])
	depositTx := &types.EncounterPartyDepositTx{
	    AccountSequence: acc.Sequence,
	    Sender:          msg.senders[1],
	    Token:           msg.Tokens[1],
	}
    rawTx = types.CDC.Marshal(depositTx)
    pubKey = acc.GetPubKey()
    isValid := pubKey.VerifySignature(rawTx, msg.signedTx)

```

- according to the signature verification result, continue follow process.

  ```js
      if (isValid) {
          // Lock assets from senders to escrow account
          escrowAccount := types.GetEscrowAddress(pool.EncounterPartyPort, pool.EncounterPartyChannel)
          // Create a deposit message
          sendMsg := banktypes.MsgSend{
  	        FromAddress: secondSenderAcc.GetAddress().String(),
  	        ToAddress:   escrowAccount.String(),
  	        Amount:      sdk.NewCoins(*msg.Tokens[1]),
          }
          k.executeTx(sendMsg)
      }
  ```

5. return execution result to ack
6. acknowledged on chain A
