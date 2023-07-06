# silver-arrow

Backend service for powering subscriptions on Lucid.

## Entities

### Account

I still need to account for how to deal with AA wallets that a created with social login.
The signer address field can be nullable and we use a different mechanism to verify subscriptions added to such accounts.

* Address (AA address)
* SignerAddress (EOA address)

### Subscription

This should include some metadata fields

* Amount
* FromAddress
* DestinationAddress
* Signature
* Interval
* Metadata (json object)

### Relationships

Address --> Subscription (one to many)

## API

RPC?
Graphql?
During account creation:
    Email, AccountAddress
Adding a new subscription/recurring payment:
    OriginAddress
    DestinationAddress
    Amount
    Token
    Interval