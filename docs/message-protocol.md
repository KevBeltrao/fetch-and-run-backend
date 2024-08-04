Format:
```
{
    "type": "playerMove",
    "payload": {
        "playerID": "player1",
        "x": 100,
        "y": 150
    }
}
```

After connecting, a initial message must be sent with the type "initial" and the playerId in the payload, like:
```
payload: {
    playerId: playerId
}
```
