# How to run
```
git clone https://github.com/voidfun/simple-blockchain.git
cd simple-blockchain/core
go run .
```

# Instructions
- Show balances: `GET {httpAddr}/balances`
- Add transactions: `POST {httpAddr}/tx/add`
  ```
  {
      "from": "3eb92807f1f91a8d4d85bc908c7f86dcddb1df57",
      "from_pwd": "security123",
      "to": "0x6fdc0d8d15ae6b4ebf45c52fd2aafbcbb19a65c8",
      "gas": 21,
      "gas_price": 1,
      "value": 10,
      "data": ""
  }
  ```
- Show the status of the node: `GET {httpAddr}/node/status`
- Sync node: `GET {httpAddr}/node/sync`
- Add peer node: `GET {httpAddr}/node/peer?ip={your_ip}&port={your_port}&miner={}&version={}`

# Features
- [x] Basic data structure
- [x] HASH/Merkle Root/Validation
- [x] Storage
- [x] Consensus(POW)
- [x] Instruction
- [x] Sync
- [x] Multi port communication
- [x] Documentation