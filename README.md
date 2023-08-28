# silver-arrow

Backend service for powering subscriptions on Lucid.

This service relies heavily on ERC-4337 primitives and zerodev [kenel](https://github.com/zerodevapp/kernel). Kernel is a modular smart contract wallet implementation that is built on the ERC-4337 standard.

ERC-4337 is an account abstraction proposal which completely avoids the need for consensus-layer protocol changes. Instead of adding new protocol features and changing the bottom-layer transaction type, this proposal introduces a higher-layer pseudo-transaction object called a UserOperation. Users send UserOperation objects into a new separate mempool. Bundlers package up a set of these objects into a single transaction by making a call to a special contract, and that transaction then gets included in a block.

By leveraging this pseudo-transaction object, user intents can be originated offchain and validated by the wallet itself onchain. This allows us to be able to implement pull payments in the smart contract wallet, inherently enabling subscriptions.

## Run locally 

- Pull the latest changes from git
- run ```go get```
- set the environment variables
- run ```APP_ENV=development go run cmd/server.go```

## Modules/Packages

### erc-4337

Contains logic to handle erc-4337 related things like communicating with the bundler, paymaster for gas abstraction, smart contract interactions and parsing user intents

### wallet

Core wallet logic, subscription handling.

### scheduler

contains logic to handle scheduling a subscription

### tests

contains integration tests

### abi

Contains abis and bindings for the smart contracts silver-arrow interacts with.

### graph

Graphql related code

### repository

The data handling layer and model lives here