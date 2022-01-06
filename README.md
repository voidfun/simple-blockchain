# How to run
```
git clone https://github.com/voidfun/simple-blockchain.git
cd simple-blockchain/core
go run .
```

# Instructions
- Show blockchain: `GET {httpAddr}/`
- Add new record to blockchain: `POST {httpAddr}/`
    ```
    {
        "BPM": 123
    }
    ```
- Register ports to blockchain: `POST {httpAddr}/nodes`
    ```
    {
        "URLs": ["http://localhost:7777", "http://localhost:8888"]
    }
    ```
- List nodes of blockchain: `GET {httpAddr}/nodes`
- Sync blocks and nodes: `GET {httpAddr}/syncAll`

# TODO list
- [x] Basic data structure
- [x] HASH/Merkle Root/Validation
- [x] Storage
- [x] Consensus(POW)
- [x] Instruction
- [x] Sync
- [x] Multi port communication
- [x] Documentation