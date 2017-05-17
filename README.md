# TRPG Chatroom

Simple server for creating HTTP chatroom that's suitable for TRPG games online.

## Installation

```
go get -u github.com/tomatorpg/tomatorpg
```

## Usage

```
chatroom
```

## Development

### Golang Server Development

The server is compatible with [gin](https://github.com/codegangsta/gin) live
reloader. You should install it for development:

```
go get -u github.com/codegangsta/gin
```

To run golang server with development mode:
```
gin
```

The live reloader will run at http://localhost:3000 by default. And it will
automatically compile your server code on file changes.


### Fronend JS Development

First you need to start a webpack-dev-server to live-reload your JS code
for development:

```
WEBPACK_DEV_SERVER_HOST=http://localhost:8080 yarn dev
```

Then you should open another terminal and start the server program. You need
to inject the webpack dev server path for development. Either:
```
WEBPACK_DEV_SERVER_HOST=http://localhost:8080 ./chatroom
```

or

```
WEBPACK_DEV_SERVER_HOST=http://localhost:8080 gin
```
