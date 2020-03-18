# CRM API test for The Agile Monkeys.

![Build status](https://davsanchez.semaphoreci.com/badges/theam_test_crm_api.svg?style=shields)

This repo contains my proposed solution to **API Test - The CRM Service** by **The Agile Monkeys**. It is written in **Go**.

## App features

### Architecture

The solution proposed makes use of a [PostgreSQL database](https://www.postgresql.org) provided by [its Docker image](https://hub.docker.com/_/postgres/), as part of a multi-container application managed with Docker Compose. The database feeds the API backend service, coded in Go, that is also enabled to run as a Docker container built from [its official image](https://hub.docker.com/_/golang/).

This configuration makes possible to setup the whole architecture with one command, provided Docker is installed, executed from the project's root directory:

```sh
docker-compose up -d --build
```

Where the option `-d` makes `docker-compose` run in detached mode and `--build` rebuilds the Go backend in case changes were made since the last setup.

To shutdown the whole architecture, a similar command needs to be used from the project's root directory:

```
docker-compose down
```

This will stop all the containers but won't remove the volumes defined for them. Thus, the database contents will be persisted in its volume for the next time it is run.

![Project architecture](./theam_test_arch.png "Project architecture")

As the above diagram suggest, it is possible to run other backends external to the Docker container architecture, provided the different configuration parameters needed (backend port, database host and ports...) are taken into account.

A stand-alone version of the backend is also running on [Heroku](https://www.heroku.com/), hooked to the GitHub repository's `heroku` branch and using the PostgreSQL database available in the SaaS platform's free tier. As long as authentication is not activated, you can make requests to the available endpoints (see [below](#API_endpoints)) at host [`https://theam-crm-api.herokuapp.com/`](https://theam-crm-api.herokuapp.com/).

The following libraries were used for the backend development:
- `gorilla/mux`: Provides the route handling via its HTTP request multiplexer.
- `lib/pq`: Go database driver for PostgreSQL.
- `crypto/bcrypt`: Cryptographic library for hashing and comparing passwords.

## <a name="API_endpoints"></a>API endpoints

The API was implemented making use of `gorilla/mux`'s router, which allow matches incoming requests against a list of registered routes and calls a handler for the route that matches the URL or other conditions. All API endpoints return a JSON object, the details below define its content for each endpoint.

For detailing the possible request inputs, conditions and outputs of the API endpoints, the following syntax is used:

```
(<condition?) <input?> -> <response>
```

For the HTTP requests, if any URL segment is a parameter, it will be enclosed in curly braces, `{variable}`, in the title and represented as `variable` in the example responses. For the example JSON responses, only descriptive values are used. No numeric value can be below `1`.

### Customers

#### `GET /customers/all`
Endpoint for getting a list of all customers in the system.
```js
(No customers) -> []
(1+ customers) -> [
    {
        "id":1,
        "name":"Customer_1_name",
        "surname":"Customer_1_surname",
        "pictureId":1,
        "lastModifiedByUserId":1
    },
    // ... (If more than 1 customer)
]
(Error) -> {"error": "error_message"}
```

<!-- #### Possible endpoint improvements
// TODO -->


#### `GET /customers/{customerId}`
Endpoint for getting the customer of a specific `customerId`.
```js
(Existing customerId) -> {
        "id":customerId,
        "name":"customer_id_name",
        "surname":"customer_id_surname",
        "pictureId":1,
        "lastModifiedByUserId":1
}
(No customers) -> {"error":"Customer not found"}
(Error) -> {"error":"error_message"}
```


#### `POST /customers/create`
Endpoint for creating a specific user in the system.
```js
(Created successfully) {
        "name":"new_customer_name",
        "surname":"new_customer_surname",
        "pictureId":1,
        "lastModifiedByUserId":1
} -> {
        "id":2,
        "name":"new_customer_name",
        "surname":"new_customer_surname",
        "pictureId":1,
        "lastModifiedByUserId":1
}
(Error) * -> {"error":"error_message"}
```

#### `PUT /customers/{customerId}`
Endpoint for updating a specific user in the system.
```js
(Updated successfully) {
        "name":"updated_customer_name",
        "surname":"updated_customer_surname",
        "pictureId":1,
        "lastModifiedByUserId":1
} -> {
        "id":customerId,
        "name":"updated_customer_name",
        "surname":"updated_customer_surname",
        "pictureId":1,
        "lastModifiedByUserId":1
}
(Nonexistent {customerId}) * -> {"error":"No customer was updated"}
(Error) * -> {"error":"error_message"}
```

#### `DELETE /customers/{customerId}`
Endpoint for deleting a specific user in the system.
```js
(Deleted successfully) * -> {"result":"success"}
(Nonexistent {customerId}) * -> {"error":"No customer was deleted"}
(Error) * -> {"error":"error_message"}
```

### Pictures

#### `GET /customers/picture/{pictureId}`
Endpoint for getting the picture path of a specific `pictureId`.
```js
(Existing pictureId) -> {
        "id":pictureId,
        "picturePath":"picture/id/path",
}
(No picture) -> {"error":"Picture not found"}
(Error) -> {"error":"error_message"}
```


#### `POST /customers/picture/upload`
Endpoint for uploading picture to the system.
```js
(Uploaded and stored successfully) [multipart_form] -> {
        "id":1,
        "picturePath":"picture/id/path.ext",
}
(Error) * -> {"error":"error_message"}
```

### User authentication and authorization

#### `POST /users/register`

#### `POST /users/login`


## Further improvements

### // TODO
// TODO

### ... and much more!
I mean, this is being done in a week, while working full time and with a world-spanning viral crisis in full force! Sure there will be room for improvement :)
