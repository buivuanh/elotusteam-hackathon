## Overview

Welcome to the Hackathon repository! This project is built with Go and PostgreSQL for data storage.
Docker and Docker Compose are used for containerization and orchestration.

## Getting Started

To get the code up and running on your local machine, follow the steps below:

### Prerequisites

- Ubuntu (18.04 or higher)
- Go (1.19 or higher)
- Docker (23.0.5 or higher)
- Docker Compose (v2.17.3 or higher)

### Installation

1. Install [Docker Engine](https://docs.docker.com/engine/install/ubuntu/)

2. Install [Docker Compose](https://docs.docker.com/compose/install/):

```shell
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

3. Install `make`: `sudo apt install build-essential`

4. Install [go](https://golang.org/dl/)

5. Clone the repository:

```shell
git clone git@github.com:buivuanh/elotusteam-hackathon.git
cd elotusteam-hackathon
```

### Build and start

- Build and start the Docker containers using Docker Compose:

```shell
make run-local
```

or

```shell
docker-compose up --build
```

or "detached" mode

```shell
docker-compose up --build -d
```

- Access the running application:
  Once the containers are up and running, you can access the application at http://localhost:8080.

- See `main_test.go` to refer some client code.

- After application is started, there are will 2 directory is created: `../tmp` (contains files is uploaded)
  and `../postgres-data` (contains database files). See more `docker-compose.yml`

- To custom env, modify `docker-compose.yml` file:

```shell
    environment:
      DATABASE_URL: <connect-string>
      JWT_SIGNING_KEY: <secret-key>
      DATA_STORE_PATH: <directory-path>
      JWT_TOKEN_EXPIRATION_HOUR: <token-expiration-hour>
```

- To run unit test:

```shell
make test-unit
```

## Project Structure

The repository is organized using the following directory structure:

- main.go: The entry point of the application.
- auth: Contains the custom middleware to authenticate and validate request.
- controller/: Contains the handlers for different HTTP routes and request handling.
- domain/: Includes the domain models and its methods.
- infrastructure/: Provides the infrastructure code, such as database connections and external service integrations.
- migrations/: Contains database migration sql scripts.
- utils/: Includes utility functions and helper code.
- test-img/: Includes some files to test upload feature.

## API Endpoints

The following API endpoints are available in this project:

### User Registration

Endpoint: `POST /register`

This endpoint allows users to register by providing a username and password. It internally hashes the password before
storing it in the database.

### User Login

Endpoint: `POST /login`

This endpoint handles user authentication by validating the provided username and password. Upon successful
authentication, a JWT token is generated and returned to the client.

### File Upload

Endpoint: `POST /upload`

This endpoint allows users to upload files, specifically images. The uploaded file must be of type "image" and should
not exceed 8 megabytes in size. The request is authenticated using JWT token verification.

## Database Schema

The project uses a PostgreSQL database with the following schema:

Table Name: users

The users table stores information about registered users, including their username and hashed password

| Column Name     | Data Type                   | Constraints      |
|-----------------|-----------------------------|------------------|
| id              | SERIAL                      | PRIMARY KEY      |
| username        | TEXT                        | NOT NULL, UNIQUE |
| hashed_password | TEXT                        | NOT NULL         |
| created_at      | TIMESTAMP without time zone | DEFAULT NOW()    |
| deleted_at      | TIMESTAMP without time zone |                  |

Table Name: images

The images table stores information about uploaded images, including file details, byte size, owner ID ...

| Column Name   | Data Type                   | Constraints                                           |
|---------------|-----------------------------|-------------------------------------------------------|
| id            | SERIAL                      | PRIMARY KEY                                           |
| file_path     | TEXT                        | NOT NULL, UNIQUE                                      |
| original_name | TEXT                        | NOT NULL                                              |
| content_type  | TEXT                        | NOT NULL                                              |
| byte_size     | INTEGER                     | NOT NULL                                              |
| owner_id      | INTEGER                     | NOT NULL, FOREIGN KEY (owner_id) REFERENCES users(id) |
| created_at    | TIMESTAMP without time zone | DEFAULT NOW()                                         |
| deleted_at    | TIMESTAMP without time zone |                                                       |