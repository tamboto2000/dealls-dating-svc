# Dealls Dating App Backend Service
Study Case for creating backend service for a dating app

## How to run
Steps to build and run the service

### 1. Build the project
Run this command to build the project into a single library
```sh
make build
```

### 2. Local stack deployment
Run this command to deploy the local stack (Postgre and Redis)
```sh
make local-deploy
```

### 3. Migration
Due to time constraints, the migration system is not finalyzed, you can manually create and seed the table by executing SQL from file `scripts/migrations/tables.sql`

### 4. Run the service
Run this command to run the service, the HTTP port is at `:8000`
```sh
./bin/dealls-dating-svc
```

## Tech stack
List of tech stack that is used for this project

### 1. PostgreSQL
As the main database to store data

### 2. Redis
As a cache for frequently accessed data

### 3. Docker
For deploying and running the service
