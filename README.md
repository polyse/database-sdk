# database-sdk

Simple sdk to add Documents in POLYSE database.

### Installing

Install PolySE Database SDK by runnig:

```bash
go get github.com/polyse/database-sdk
```

### Usage 

1) Import package `import sdk "github.com/polyse/database-sdk"` 
2) Start [polySE database](https://github.com/polyse/database) on _<example_host>:<example_port>_
3) Create new client like :
      ``` go
         newclient, err := sdk.NewDBClient("<example_host>:<example_port>")
      ```
4) Use client end enjoy:).