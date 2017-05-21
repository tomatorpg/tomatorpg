# TomatoTRPG [![Build Status][travis-shield]][travis-link]

Simple server for creating HTTP chatroom that's suitable for TRPG games online.

[travis-link]: https://travis-ci.org/tomatorpg/tomatorpg
[travis-shield]: https://api.travis-ci.org/tomatorpg/tomatorpg.svg?branch=master

## Installation

```
go get -u github.com/tomatorpg/tomatorpg/cmd/tomatorpg
```

The server code will be inside your [`$GOPATH`](https://golang.org/doc/code.html#GOPATH)
folder in `$GOPATH/src/github.com/tomatorpg/tomatorpg`.

## Usage

If you have included `$GOPATH/bin` into your `$PATH`, you can run this:

```
tomatorpg
```

Or you can run this:

```
$(go env GOPATH)/bin/tomatorpg
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
gin -port 8080 -build cmd/tomatorpg
```

The live reloader will run at http://localhost:8080 by default. And it will
automatically compile your server code on file changes.


### Fronend JS Development

First you need to start a webpack-dev-server to live-reload your JS code
for development:

```
NODE_ENV=development yarn dev
```

Then you should open another terminal and start the server program. You need
to inject the webpack dev server path for development. Either:
```
NODE_ENV=development gin
```

or if you have no need to change the go code, you may test against server
binary with:

```
NODE_ENV=development ./tomatorpg
```

### `.env` file

Both the Golang server and the webpack-dev-server config support loading
variables from dotenv file (`.env`).  You may simply put your environment
variables in the `.env` of the project root.

For production, you may add a `.env` file at the same folder of the `tomatorpg`
binary.
