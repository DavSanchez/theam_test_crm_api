# CRM API test for The Agile Monkeys.

This repo contains my proposed solution to API Test - The CRM Service by The Agile Monkeys. It is written in Go.

## App features

### Architecture

The solution proposed makes use of a [PostgreSQL database](https://www.postgresql.org) provided by [its Docker image](https://hub.docker.com/_/postgres/), as part of a multi-container application managed with Docker Compose. The database feeds the API backend service, coded in Go, that is also enabled to run as a Docker container built from [its official image](https://hub.docker.com/_/golang/).

This configuration makes possible to setup the whole architecture with one command, executed from the project's root directory:

```sh
docker-compose up -d --build
```

Where the option `-d` makes `docker-compose` run in detached mode and `--build` rebuilds the Go backend in case changes were made since the last setup.

To shutdown the whole architecture, a similar command needs to be used from the project's root directory:

```
docker-compose down
```

This will stop all the containers but won't remove the volumes defined for them. Thus, the database contents will be persisted in its volume for the next time it is run.

![Project architecture](./be_architecture.png "Project architecture")

As the above diagram suggest, it is possible to run other backends external to the Docker container architecture, provided the different configuration parameters needed (backend port, database host and ports...) are taken into account.

The following libraries were used for the backend development:
- `gorilla/mux`: Provides the route handling via its HTTP request multiplexer.
- `lib/pq`: Go database driver for PostgreSQL.
- `crypto/bcrypt`: Cryptographic library for hashing and comparing passwords.

## API endpoints

The API was implemented making use of `gorilla/mux`'s router, which allow matches incoming requests against a list of registered routes and calls a handler for the route that matches the URL or other conditions. All API endpoints return a JSON object, the details below define its content for each endpoint.

### `GET /customers/all`
// TODO

#### Possible endpoint improvements
// TODO


### `GET /customers/{customerId}`
// TODO


### `POST /customers/create`
// TODO


### `PUT /customers/{customerId}`
// TODO


### `DELETE /customers/{customerId}`
// TODO


## Further improvements

### // TODO
// TODO

### ... and much more!
I mean, this was done in a week, while working full time and with a world-spanning viral crisis in full force! Sure there is a lot to troubleshoot :)
