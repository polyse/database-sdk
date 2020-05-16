# database-sdk

Simple sdk to add Documents in POLYSE database.

#### Usage :

1) Start POLYSE database on _localhost:8080_
2) Create collections, for example _news_
3) Create new client like :
    > newclient := database_sdk.NewDBClient("localhost:8080")
4) Send data to save your new documents like :
    > docs, err := newclient.SaveData("news", []Documents{...})
5) Handle error.