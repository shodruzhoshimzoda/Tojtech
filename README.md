# Tojtech API

A robust Go-based REST API backend for an e-commerce platform specializing in IT hardware and computer products.

First of all ensure that you have this stuff on your system
1) Go (v.1.22) or latest
2) PostgreSQL

# Database setup
Before running the application, make sure PostgreSQL is running and create the target database:
```
    CREATE DATABASE myshop;
```

# Installation and setup
```
    git clone https://github.com/shodruzxoshimzoda/tojtech.git
    cd tojtech
```

# Download dependencies

```
    go mod tidy
```

# Install the migration tool, if you have not:
    
```
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
# Apply database migrations:
Set your database password in the current shell session and run the migrations via Makefile

```
 DB_PASSWORD='your_postgres_password'
make migrate-up

```

# RUN SERVER

You can run server with two methods

1) go run ./cmd/api/ -db-password=sabrr1se
2) DB_PASSWORD=your_db_password make server-run
