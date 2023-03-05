# Node & Go gRPC demo using buf.build

This codebase hosts a client and server gRPC service showcasing buf.build as well as using packages which are hosting your proto files and generated code.


## Intro

This repo using [proto-store](https://github.com/jtoloui/proto-store) package in both it's go and javascript example.

The servers all connect to a mongo db which you can connect your own by setting in an .env file at all entry levels within the go/js server folders.

The model in the db collection used in the example can we see [here](./javascript/server-grpc/models/movie.ts) or below

```ts
const movieSchema = new Schema({
	title: {
		type: String,
		required: true,
	},
	year: {
		type: Number,
		required: true,
	},
	director: {
		type: String,
		required: true,
	},
});
```

The endpoints in both languages provide examples of

Port: `:8080`
|Method type|Endpoint|Params/Body|
|--|--|--|
|GET|/|Params: page, per_page|
|GET|/:id|N/A|
|POST|/|Body: title, director, year|

## Go

The client for the go example uses Gin

## Javascript

The Javascript using typescript as for the client server it's running an express server


## Local development

From the root director you can run the commands inside of the [Makefile](./Makefile)

### Javascript
e.g running javascript client and service

```bash
make js-server
```

```bash
make js-client
```

### Go
e.g running go client and service

```bash
make go-run-grpc
```

```bash
make go-run-gin
```