# What's the diffiences between the HTLC and IBC Atomic Swap

## HTLC Concept (from Irisnet)

Hash Time Locked Contract (HTLC) (opens new window)is a protocol intended for cross-chain atomic swaps. It performs a P2P swap in a decentralized form. HTLC can guarantee that all swap processes among participating parties take place or nothing is done. Leveraging the HTLC on IRIS Hub, users can implement cross-chain asset swaps with Bitcoin, Ethereum, LTC and all other chains which support HTLC

## Steps of HLTC:

Supposed that Alice swaps BTC for IRIS with Bob. The atomic swap process between Alice and Bob can be divided into following steps.

- step 0
Alice and Bob reach an agreement on the swap by an off-chain procedure, which is completed commonly on an exchange platform. Let Alice be a maker and Bob be a taker. The agreement includes the exchange rate between BTC and IRIS, the each other's address on the counterparty chain and a unique hash lock generated from a secret, and an appropriate time span.

- step 1
Afterwards, the secret holder, Bob sends an HTLC transaction which commits to transfer IRIS of the specified amount to Alice(the maker) with the negotiated hash lock on IRIS Hub.

- step 2
An event is emitted on IRIS Hub which indicates that Bob has created an HTLC. Alice is informed of this event by monitor tools (usually wallets) or the platform. The HTLC initiating transaction with the same hash lock will be sent to Bitcoin by Alice once the event is validated against the agreement. Particularly the HTLC will be locked by an quite smaller time span than the one provided by Bob.

- step 3
Bob is informed of and confirms the event on Bitcoin. Then Bob claims the HTLC-locked BTC with the owned secret before the time span set by Alice on Bitcoin.

- step 4
The secret will be disclosed while the HTLC is claimed successfully by Bob on Bitcoin. Alice will perfom the same claim to the locked IRIS with the secret before the expiration time on IRIS Hub.

## IBC Atomic Swap (IAS)

## Steps of IAS

- step 0 


