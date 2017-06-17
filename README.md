# TomatoRPG [![Build Status][travis-shield]][travis-link]

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

You may work in the [Frontend Development mode](#frontend-development). Or
you may work in the [Fullstack Development mode](#fullstack-development).

Depends on what development work are you doing.


### Fronend Development

To golang server on development mode, with webpack-dev-server for frontend
development:

```
yarn devfront
```

The frontend javascript / scss code will be monitored and rebuilt by
webpack-dev-server. The golang server will be using the javascript that
provided by webpack-dev-server.


### Fullstack Development

You should install [gin][gin] live reloader for fullstack development:

```
go get -u github.com/codegangsta/gin
```

To run development server both for frontend and backend:

```
yarn dev
```

The golang server binary will be monitored, automatically rebuilt by
[gin][gin] and listen at http://localhost:8080 by default.

The frontend javascript / scss code will be monitored and rebuilt by
[webpack-dev-server][webpack-dev-server]. The javascript built assets
will be served at http://localhost:8081/assets by default.

[gin]: https://github.com/codegangsta/gin
[webpack-dev-server]: https://www.npmjs.com/package/webpack-dev-server


### Resolve port conflict

If you encountered port conflict, you may use the following env variables:

* `PORT` controls the port for tomatorpg server to bind to (if you're not
  running the server through [gin][gin]). Default is `8080`.

* `PUBLIC_URL` provides the base public URL of the site to OAuth2. This will
  be used in the OAuth2 login (i.e. Google / Facebook / Twitter) process only.
  Default is: `http://localhost:8080`.

* `WEBPACK_DEV_SERVER_HOST` defines the webpack-dev-server public path. If
  you have set `NODE_ENV` to "`development`", the tomatorpg server will use
  javascript from this host, instead of its built-in assets.

You may use the `.env` file to override these variables. Or you may directly
set these variables in your environment.


### `.env` file

Both the Golang server and the webpack-dev-server config support loading
variables from dotenv file (`.env`).  You may simply put your environment
variables in the `.env` of the project root.

For production, you may add a `.env` file at the same folder of the `tomatorpg`
binary.
