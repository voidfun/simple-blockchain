# How to run
```
git clone https://github.com/voidfun/simple-blockchain.git
cd simple-blockchain/core
go run .
```

# Instructions
- Show balances: `GET {httpAddr}/balances/list`
- Add transactions: `GET {httpAddr}/transactions: `POST {httpAddr}/tx/add`
- Show the status of the node: `GET {httpAddr}/node/status`
- Sync node: `GET {httpAddr}/node/sync`
- Add peer node: `POST {httpAddr}/node/peer`

# Features
- [x] Basic data structure
- [x] HASH/Merkle Root/Validation
- [x] Storage
- [x] Consensus(POW)
- [x] Instruction
- [x] Sync
- [x] Multi port communication
- [x] Documentation