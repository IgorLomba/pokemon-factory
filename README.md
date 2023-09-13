# Pokemon Factory

Welcome to the Pokemon Factory project! This project provides a set of APIs to generate and manage Pokemons
(please Nintendo don't sue me).

Created for study purposes, this project is a simple example of a REST API using Go, Docker, Swagger.
Containing Fifo and Lifo queues to generate and get pokemons.

## Getting Started

### Docker

Run:

```bash
docker compose up -d
```

Clean:

```bash
docker compose down --volumes --remove-orphans --rmi all
```

```bash
docker compose rm --force --stop --volumes
```

### API Documentation

Swagger:

```bash
swagger generate spec -m -o swagger.json && swagger serve -F swagger swagger.json
``` 

### Thanks!
