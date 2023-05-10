# LilShorty

A simple URL shortener built with Go and MongoDB.

## Requirements

- Go (v1.16 or higher)
- MongoDB

## Installation

1. Clone the repository:
```shell
git clone https://github.com/cchhaarroonn/LilShorty.git
```
2. Change to the project directory:
```shell
cd LilShorty
```
3. Initialize project and get dependencies:
```shell
go mod init <charon/lilshorty>
go get github.com/gin-gonic/gin
go get go.mongodb.org/mongo-driver/mongo
```
4. Set up MongoDB by creating a database and a collection. Update the MongoDB connection URL in the code:
```go
client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://<username>:<password>@<host>:<port>")) <line 22>
```

## Usage

1. Run the application:
```shell
go run main.go
```
2. The server will start running on locahost:42069 by default change it on line 127

## Shortening a URL
Send a POST request to /createShort/:url to shorten a URL. Replace :url with the actual URL you want to shorten.

```shell
POST /createShort/http://example.com
```

```json
{
  "status": "URL shortened, key is: <generated-key>"
}
```

## Redirecting to the Original URL

```shell
GET /short/<generated-key>
```

```json
{
  "status": "<original-url>"
}
```

## Contributing
Contributions are welcome! If you find any issues or want to add new features, please open an issue or submit a pull request.

## Note

As always there is a note lmao ... I might make website for this like basic html, css, js but it really depends if I want to burn my brain at some point until then just enjoy using this simple thingy it is not much but is quite good
