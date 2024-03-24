# Money Transfer Backend Service
This is a money transfer that allows users to make withdraw and deposit.

# Installation
To get started with this project you must install docker and docker-compose. You can find the installation guide [here](https://docs.docker.com/get-docker/).
```bash
docker-compose up
```
# Documents
## API
To get api documentation you can visit [http://localhost:8080/docs](http://localhost:8080/docs), or you can find in path `./api/docs/`.
## Structure
```
.
├── Dockerfile
├── Makefile
├── README.md
├── api                         # define api for service
│   ├── docs                    # api documentation generated from proto file
│   └── proto                   # define proto for grpc
├── cmd                         # command to execute service: http, cronjob, worker, ....
├── configs                     # config for service
│   └── http_server             # config for http server
├── docker-compose.yml
├── internal                    # contains all the business logic of the application
│   ├── adapters                # contains the implementation of the application ports
│   ├── applications            # business logic of the application
│   │   ├── accounts            # business logic of accounts
│   │   ├── transactions        # business logic of transactions withdraw and deposit
│   │   └── users               # business logic of users
│   ├── entities                # contains the domain entities
│   │   ├── account             # account entity
│   │   ├── transaction         # transaction entity
│   │   │   ├── deposit.go      # deposit transaction
│   │   │   ├── transaction.go  # transaction entity contains withdraw and deposit transaction
│   │   │   └── withdraw.go     # withdraw transaction
│   │   └── user
│   └── ports                   # contains the application ports, expose the application to the outside world: http, grpc, cronjob, worker, ...
│       ├── grpc                # grpc servers
│       │   ├── accounts        # grpc accounts server
│       │   ├── transactions    # grpc transactions server
│       │   └── users           # grpc users server
│       └── http                # http servers
└── main.go                     # entry point of the application
```

# Requirements
- Implement a Rest API with CRUD functionality.
- Database: MySQL or PostgreSQL.
- Unit test as much as you can.
- Set up service with docker-compose.
- Secure the API with your choice of authentication mechanism.
### Data models
#### User
- has accounts
### Account
- has at least the following fields:
- `name`: name of account
- `bank`: name of bank (3 possible values: `VCB`, `ACB`, `VIB`)
- has transactions
### Transaction
- has at least the following fields:
- `amount`: amount of money
- `transaction_type`: type of transaction (2 possible values: `withdraw`, `deposit`)

# Database
MySQL as database, use sqlc for generating code from sql file.
Schema of database in folder `./internal/adapters/repository/sqlc/schema`.

## Testing

## Deployment
To deploy this project you can use docker-compose, kubernetes, or any other container orchestration tools.
### Docker-compose
```bash
make docker-compose-up
```