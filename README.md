# Client-Server API Go

This project consists of two Go programs (`client.go` and `server.go`) that demonstrate how to consume an external API, save data in a SQLite database, and manipulate files using standard Go packages and external packages.

## Description

- **client/client.go**: A client that makes an HTTP request to the server (`server/server.go`) to request the USD to BRL exchange rate and saves the received value in a text file.
- **server/server.go**: A server that consumes an external API to obtain the USD to BRL exchange rate, returns the result to the client in JSON format, and logs the exchange rate in a SQLite database.

## Requirements

- Go 1.18 or higher
- SQLite3

## Setup

1. **Initialize Go Modules**:
    ```sh
    go mod init github.com/pansani/client-server-api
    ```

2. **Install Dependencies**:
    ```sh
    go get github.com/mattn/go-sqlite3
    ```

3. **Structure**:
    Ensure your project structure looks like this:
    ```
    client-server-api-go/
    ├── client/
    │   └── client.go
    ├── server/
    │   └── server.go
    ├── cotacao.txt
    ├── cotacoes.db
    ├── go.mod
    └── go.sum
    ```

## Running the Project

### Running the Server

1. Open a terminal and navigate to the `server` directory.
2. Run the server:
    ```sh
    go run server.go
    ```
3. The server will start and listen on port 8080.

### Running the Client

1. Open another terminal and navigate to the `client` directory.
2. Run the client:
    ```sh
    go run client.go
    ```
3. The client will make a request to the server, receive the exchange rate, and save it in the `cotacao.txt` file.

## Functionality

- **Server**:
  - Fetches the USD to BRL exchange rate from the API `https://economia.awesomeapi.com.br/json/last/USD-BRL` with a timeout of 200ms.
  - Saves the fetched exchange rate in a SQLite database with a timeout of 10ms.
  - Provides an endpoint `/cotacao` that returns the current exchange rate in JSON format.
  
- **Client**:
  - Makes a request to the server's `/cotacao` endpoint with a timeout of 300ms.
  - Saves the received exchange rate in a file named `cotacao.txt` in the format: `Dólar: {value}`.

## Error Handling

- The server logs errors if the request to the external API or the database operation exceeds the timeout.
- The client logs errors if the request to the server exceeds the timeout or if there is an error saving to the file.

## Example Output

- **cotacao.txt**:
    ```
    Dólar: 5.42
    ```
