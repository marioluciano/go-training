# Go Training

This repository contains three small Go challenges:

- [client-server-api](client-server-api)
- [hello-world](hello-world)
- [multitheading](multitheading)

## 1. Hello World

This example prints a simple message.

From the repository root:

```bash
go run ./hello-world
```

Or from the challenge folder:

```bash
cd hello-world
go run hello.go
```

## 2. Client/Server API

This challenge has two parts:

1. Start the server so it exposes the quote endpoint on port 8080.
2. Run the client so it fetches the quote and saves it to a file.

Open one terminal and start the server:

```bash
go run ./client-server-api/server
```

Or from the server folder:

```bash
cd client-server-api/server
go run server.go
```

Open a second terminal and run the client:

```bash
go run ./client-server-api/client
```

Or from the client folder:

```bash
cd client-server-api/client
go run client.go
```

The client writes the fetched quote to a file named cotacao.txt in the current working directory.

## 3. Multithreading

This example makes concurrent requests to two different CEP APIs and prints the first successful response.

From the repository root:

```bash
go run ./multitheading
```

Or from the challenge folder:

```bash
cd multitheading
go run main.go
```

## Notes

- The module root is already configured, so you can run the examples directly with Go.
- The client/server example depends on the server being running before the client starts.
- The folder for the third challenge is currently named [multitheading](multitheading), which is the name used in this repository.