# TomatoRPG Protocol

## Basics

TomatoRPG is using websocket as the transport layer for RPC calls from web
client to the server.

## RPC Request

RPC protocol is based on JSON.

Either `endpoint` must be present for a valid API call.

```json
{
  "tomatorpc": "0.1",
  "context": "session",
  "entity": "roomActivities",
  "action": "create",
  "payload": {
    "action": "some activitiy",
    "name": "some name"
  }
}
```

Attribute `type` is reserved for redux compatibility and will be ignored by
server.

## Server Message Stream

Server send 2 types of messages to clients:

1. Activity broadcast;
2. RPC response.

Activity broadcast:

```json
{
  "tomatorpc": "0.1",
  "entity": "roomActivities",
  "type": "broadcast",
  "data": {
    "room_id": "room-id",
    "user_id": "id of the acting user",
    "character_id": "id of the acting character, if applicable",
    "action": "room activity type",
    "message": "some message",
    "timestamp": "2016xxxx"
  }
}
```

RPC response
```json
{
  "tomatorpc": "0.1",
  "entity": "related entity",
  "type": "response",
  "data": {
    "...": "..."
  }
}
```
