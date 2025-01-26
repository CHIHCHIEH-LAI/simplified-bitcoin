# Simplified BitCoin
## Requirements
### Functional
- Allow users to mine blocks
- Support transactions
- Store the blockchain consistently
- Synchronize data across nodes

### Non-Functional
TBD

## Execute
### Start a Long-Running P2P Node on Simplified-Bitcoin Blockchain

#### Start the first ever node
```
go run cmd/node/main.go -port=8080 -address=127.0.0.1:8080
```

Explanation of Flags
- -port: The port on which the node will listen for incoming connections (e.g., 8080).
- -address: The IP address and port of the current node (e.g., 127.0.0.1:8080).

#### Start a node that joins an existing P2P network and connects to the bootstrap node
```
go run cmd/node/main.go -port=8081 -address=127.0.0.1:8081 -bootstrap=127.0.0.1:8080
```

Explanation of Flags
- -port: The port on which the node will listen for incoming connections (e.g., 8081).
- -address: The IP address and port of the current node (e.g., 127.0.0.1:8081).
- -bootstrap (optional): The address of a bootstrap node to join the existing P2P network (e.g., 127.0.0.1:8080).

### Create a Wallet with a Private Key and a Public Key
```
go run cmd/wallet/main.go -action=createWallet -wallet=wallet.json
```

Explanation of Flags
- -action: Action to perform
- -walletFile: The filename for saving the wallet

### Create a Transaction
```
go run cmd/wallet/main.go -address=127.0.0.1:8081 -bootstrap=127.0.0.1:8000 -action=createTx -wallet=wallet.json -recipient=user1 -amount=0.01 -fee=0.001
```

## Improvements
