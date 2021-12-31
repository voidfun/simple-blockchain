# How to run
```
git clone https://github.com/voidfun/simple-blockchain.git
cd simple-blockchain
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
- Register ports to blockchain: `POST {httpAddr}/nodes/register`
    ```
    {
        "Ports": ["9001", "9002"]
    }
    ```
- List nodes of blockchain: `GET {httpAddr}/nodes`

# TODO list
- [ ] Basic data structure
- [ ] HASH/Merkle Root/Validation
- [ ] Storage
- [ ] Consensus(POW)
- [ ] Instruction
- [ ] Sync
- [ ] Multi port communication
- [ ] Documentation