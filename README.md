# Go CRUD api

Small Golang project for an interview

The goal is to develop a simple CRUD API using go and its `net/http` library.

## Usage

```
go run main.go
```

You can then test it with `httpie`.

```
http GET "localhost:8080/hello?lang=fr"
http GET "localhost:8080/hello?lang=it"
http POST "localhost:8080/hello" "hello=hola" "language=it"
http GET "localhost:8080/hello?lang=it"
http DELETE "localhost:8080/hello?lang=it"
http GET "localhost:8080/hello?lang=it"
```

You can specify the host, port and provider with args.

```
go run main.go -host localhost -port 8080 -provider map
```

Current greetings *providers* are:
- map: in memory
- redis: specify the host with "-redisHost"

## Tests

```
go test
```

## Docker

A dockerfile is available for deploying on a distributed environment.

```
#Define you custom port
PORT=2000
docker build -t dixneuf19/go-crud-api .
docker run -p "$PORT:8080" dixneuf19/go-crud-api
```


## What's next ?

- Logging
- Docker-compose