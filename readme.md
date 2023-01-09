# Cryptocurrency websocket gateway

### You may subscribe to real-time changes on any available table.

Endpoints:
WS `ws://0.0.0.0/ws` - websocket connection uri

Basic commands:

- subscribe:

```json
{
  "action": "subscribe",
  "symbols": [
    "string",
    "string"
  ]
}
```

- unsubscribe:

```json
{
  "action": "unsubscribe"
}
```

GET `http://0.0.0.0:80/statistics` - get statistics by cryptocurrency symbol

Params:

* limit - int, limit of data in response, min=1, max=1000
* last_id - int, last id from response for pagination
* symbol - string, cryptocurrency symbol filter

Response:

```json
{
  "msg": "Success",
  "data": [
    {
      "id": 170,
      "code": "XBTUSDTH23",
      "value": 17123.32,
      "date": "2023-01-08T17:48:23Z"
    },
    {
      "id": 274,
      "code": "XBTUSDTH23",
      "value": 17123.01,
      "date": "2023-01-08T17:48:23Z"
    }
  ]
}
```