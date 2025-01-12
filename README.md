# Simplified BitCoin
## Requirements
### Functional
- Allow users to mine blocks
- Support transactions
- Store the blockchain consistently
- Synchronize data across nodes

### Non-Functional
TBD

## Core Components
- Command-Line Interface (CLI): User-facing component for node control
    - Parse user commands and trigger actions
    - Communicate with other components (e.g. start mining, create transactions)
- Blockchain Core: Manages blocks, transactions, and validation
    - Define block, blockchain, and transaction structures
    - Handle adding blocks, chain validation, and synchronization
- Database Layer: Handles persistent storage
    - Persist blocks and metadata locally
    - Provide retrieval methods
- Mining: Responsible for proof-of-work and block creation
    - Provide proof-of-work to find valid hashes
    - Gather transactions from the mempool
- Networking (P2P): Manages communication between nodes
    - Broacast transactions and blocks to peers
    - Synchronize with the longest chain
- Crypto Utilities: Provides hashg and signing utilities
    - Hashing for proof-of-work and block integrity
    - Signing/verification for transaction authenticity

## Plan for Extensibility
- Add wallet or private keys
- Add advanced consensus mechanisms (e.g. proof-of-stake)

