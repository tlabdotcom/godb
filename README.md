# GoDB Package Documentation

## Environment Configuration

### Redis Configuration

To configure the Redis client, you need to set the following environment variables:

- `REDIS_HOST`: The host and port of the Redis server. Default is `localhost:6379`.
- `REDIS_PASSWORD`: The password for the Redis server (if applicable).
- `REDIS_INDEX_DB`: The index of the Redis database to use (default is 0).
- `REDIS_TIMEOUT`: (Optional) The timeout duration for connecting to Redis. Default is 10 seconds.
- `REDIS_POOL_SIZE`: (Optional) The maximum number of connections in the Redis connection pool. Default is 10.
- `REDIS_MAX_RETRIES`: (Optional) The maximum number of retries for Redis commands. Default is 3.

### PostgreSQL Configuration

To configure the PostgreSQL connection, set the following environment variables:

- `DB_POSTGRESQL_DSN`: The Data Source Name (DSN) for the PostgreSQL database connection. It should include the username, password, host, port, and database name (e.g., `user:password@tcp(localhost:5432)/dbname`).
- `MAX_OPEN_CONNS`: (Optional) The maximum number of open connections to the database. Default is 20.
- `MAX_IDLE_CONNS`: (Optional) The maximum number of idle connections to the database. Default is 10.
- `CONN_MAX_LIFETIME`: (Optional) The maximum lifetime of a connection. Default is 1 hour.
- `ENABLE_QUERY_DEBUG`: (Optional) Set to `true` to enable query debugging.

## Package Usage

### Redis Package
```go
client := godb.GetRedis()
```

### PostgreSQL Package
```go
db := godb.GetPostgresDB()

```