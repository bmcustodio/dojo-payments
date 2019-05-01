# dojo-payments

## Running

The API server requires access to a MongoDB database in order to store data.
By default, it tries to connect to MongoDB at `mongodb://localhost:27017`, and to use a database called `dojo-payments`.
The easiest way to get a compatible MongoDB setup for testing purposes is to use Docker:

```shell
$ docker run --detach --name dojo-payments-mongodb --publish 27017:27017 mongo:4.0.9
```  

To run the API server, you may then run

```shell
$ make run
```

This command starts the API server at `http://localhost:8080`.
In case you want the API server to serve requests at a different host or port, you must instead run

```shell
$ make run BIND_ADDR="<host>:<port>"
```

replacing `<host>` and `<port>` with the desired host and port.
Likewise, in case you want the API server to connect to MongoDB at a different URL or to use a different database, you must instead run

```shell
$ make run MONGODB_URL="<mongodb-url>" MONGODB_DATABASE="<mongodb-database>"
```

replacing `<mongodb-url>` and `<mongodb-database>` with the desired values.

## Testing

In order to run the end-to-end test suite, you may run

```shell
$ make test.e2e
```

In case the API server is not serving requests at `http://localhost:8080`, you must instead run

```shell
$ make test.e2e BASE_URL="http://<host>:<port>"
```

replacing `<host>` and `<port>` with the host and port where the API server can be reached.         

## Payments API

### Creating a payment

To create a payment, you may run

```shell
$ curl -X POST http://localhost:8080/payments \
  -H 'Content-Type: application/json' \
  -d '{
        "beneficiary": {
          "account_number": "1234",
          "bank_id": "4321",
          "name": "John"
         },
        "debtor": {
          "account_number": "5678",
          "bank_id": "8765",
          "name": "Dave"
        },
        "amount": 314.15,
        "currency": "EUR",
        "date": "2019-04-30T22:30:00Z",
        "description": "Order #1"
      }'
```

### Listing payments

To list all payments, you may run

```shell
$ curl -X GET http://localhost:8080/payments
```

### Getting a payment by ID

To get a payment by its ID (e.g. `5cc9ba4ee3e758d97d491b6a`), you may run

```shell
$ curl -X GET http://localhost:8080/payments/5cc9ba4ee3e758d97d491b6a
```

### Updating a payment by ID

To update a payment by its ID (e.g. `5cc9ba4ee3e758d97d491b6a`), you may run

```shell
$ curl -X PUT http://localhost:8080/payments/5cc9ba4ee3e758d97d491b6a \
  -H 'Content-Type: application/json' \
  -d '{
        "beneficiary": {
          "account_number": "1234",
          "bank_id": "4321",
          "name": "John"
        },
        "debtor": {
          "account_number": "5678",
          "bank_id": "8765",
          "name": "Dave"
        },
        "amount": 314.15,
        "currency": "USD",
        "date": "2019-04-30T22:30:00Z",
        "description": "Order #1 (Fixed)"
      }'
```

### Deleting a payment by ID

To delete a payment by its ID (e.g. `5cc9ba4ee3e758d97d491b6a`), you may run

```shell
$ curl -X DELETE http://localhost:8080/payments/5cc9ba4ee3e758d97d491b6a
```

## License

Copyright 2019 Bruno Miguel Custodio

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
