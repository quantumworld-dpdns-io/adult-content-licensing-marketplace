# Local Postgres Flow

## Start dependencies

```bash
make docker-up
```

## Set environment

```bash
export LICENSE_STORE_BACKEND=postgres
export DATABASE_URL='postgres://postgres:postgres@localhost:5432/marketplace?sslmode=disable'
export AUTH_SECRET='dev-secret'
```

## Bootstrap schema

```bash
make migrate
```

## Run service

```bash
make run
```

## Run integration test

```bash
INTEGRATION_DB=1 DATABASE_URL='postgres://postgres:postgres@localhost:5432/marketplace?sslmode=disable' go test ./internal/license -run TestSQLRepositoryCreateListGet -v
```

## Run audit store integration test

```bash
INTEGRATION_DB=1 DATABASE_URL='postgres://postgres:postgres@localhost:5432/marketplace?sslmode=disable' go test ./internal/audit -run TestSQLStoreAppendList -v
```
