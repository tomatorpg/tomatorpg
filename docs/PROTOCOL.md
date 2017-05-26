# TomatoRPG Protocol

## Basics

TomatoRPG is using websocket as the transport layer for RPC calls from web
client to the server.

## RPC Request

RPC protocol is based on JSON.

RPC Request is divided into 2 groups: CRUD and PubSub.

### CRUD Requests

CRUD is the main way to list, create, remove, update or delete entities on the
server. The field "`action`" is by default "`create`". You may omit the
parameter.

Please note that `id` should be unique between requests of the same websocket
session. It serves as the identifier to identify response for that request.

```json
{
  "tomatorpc": "0.2",
  "id": "some-client-unique-id",
  "group": "crud",
  "entity": "rooms",
  "method": "create",
  "params": {
    "name": "room name",
    "...": "..."
  },
}
```

```json
{
  "tomatorpc": "0.2",
  "id": "some-client-unique-id",
  "group": "crud",
  "entity": "roomActivities",
  "method": "create",
  "params": {
    "room_id": "room-id",
    "user_id": "user-id",
    "char_id": "char-id",
    "type": "message",
    "message": "This is the moment of truth",
    "...": "..."
  },
}
```

```json
{
  "tomatorpc": "0.2",
  "id": "some-client-unique-id",
  "group": "crud",
  "entity": "roomActivities",
  "method": "create",
  "params": {
    "room_id": "room-id",
    "user_id": "user-id",
    "char_id": "char-id",
    "type": "dice",
    "message": "1d6 > 3",
    "...": "..."
  },
}
```

### PubSub Request

PubSub requests are the ones related to the websocket session. Such as subscribe
to / unsubscribe to room activities.

```json
{
  "tomatorpc": "0.2",
  "id": "some-client-unique-id",
  "group": "pubsub",
  "method": "subscribe",
  "data": {
    "id": "room-id"
  }
}
```

```json
{
  "tomatorpc": "0.2",
  "id": "some-client-unique-id",
  "group": "pubsub",
  "method": "unsubscribe",
  "data": {
    "id": "room-id"
  }
}
```

```json
{
  "tomatorpc": "0.2",
  "id": "some-client-unique-id",
  "group": "pubsub",
  "method": "whoami",
}
```

```json
{
  "tomatorpc": "0.2",
  "id": "some-client-unique-id",
  "group": "pubsub",
  "method": "whereami",
}
```


Attribute `type` is reserved for redux compatibility and will be ignored by
server.

## Server Message Stream

Server send 2 types of messages to clients:

1. RPC response
2. Broadcast

### RPC response

```json
{
  "tomatorpc": "0.2",
  "id": "some-client-unique-id",
  "message_type": "response",
  "status": "success",
  "group": "crud",
  "entity": "roomActivities",
  "method": "create",
}
```

```json
{
  "tomatorpc": "0.2",
  "id": "some-client-unique-id",
  "message_type": "response",
  "group": "pubsub",
  "method": "whoami",
  "status": "success",
  "data": {
    "id": "user-id",
    "name": "name of user",
  }
}
```

or, if failed

```json
{
  "tomatorpc": "0.2",
  "id": "some-client-unique-id",
  "message_type": "response",
  "group": "crud",
  "entity": "roomActivities",
  "method": "create",
  "status": "error",
  "error": "some error messages",
}
```


### Activity broadcast

```json
{
  "tomatorpc": "0.2",
  "id": "some-client-unique-id",
  "message_type": "broadcast",
  "data": {
    "room_id": "room-id",
    "user_id": "id of the acting user",
    "char_id": "id of the acting character, if applicable",
    "type": "dice",
    "message": "2d6 > 7",
    "dice": "2 3",
    "dice_result": "failed",
    "timestamp": "2016xxxx"
  }
}
```
