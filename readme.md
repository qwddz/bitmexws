# Cryptocurrency websocket gateway

### You may subscribe to real-time changes on any available table.

Endpoints:
`ws://0.0.0.0/ws` - websocket connection uri

Basic commands:

- subscribe:

```json
{
  "action": "subscribe",
  "symbols": [
    "string", "string"
  ]
}
```

- unsubscribe:
```json
{
  "action": "unsubscribe"
}
```