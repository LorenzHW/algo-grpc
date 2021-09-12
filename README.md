A gRPC server to handle basic functionality of the Algorand SDK:
- Create a wallet and one account that belongs to that wallet
- Get the information of a specific account
- Make a transaction between two accounts
- Get all transactions for a specific account
- Get one block on the chain, and the information associated with that block

### Prerequisite
- [Algorand Sandbox](https://github.com/algorand/sandbox)
- [Algorand Go SDK](https://github.com/algorand/go-algorand-sdk/)
- [Go gRPC](https://github.com/grpc/grpc-go)


### Run the server
Make sure that the Algorand Sandbox is running.

Start the server:
```
go run ./server/.
```

### Make a request
```
go run ./client/.
```
