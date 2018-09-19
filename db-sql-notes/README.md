## https://golang.org/pkg/database/sql/

### `DB` type

DB is a database handle representing a pool of zero or more underlying connections. It's safe for concurrent use by multiple goroutines.

**You can use this type to create statements and transactions, execute queries, and get results**

> sql.DB is not a database connection

**sql.DB performs some important tasks behind the scenes:**
• Open and close the connection to the actual **underlying database through the driver.**
• It manages a pool of connections as needed, which may be a variety of things.

sql.DB abstraction is designed so that you don't have to worry about managing concurrent access to the underlying data store. A connection is marked as available when it is used to perform a task, and then returns to the available pool when it is not in use. One of the consequences of this is that if you can't release the connection to the pool, it can cause db.SQL to open a lot of connections, which can run out of resources

To use database/sql , you need database/sql itself, as well as the specific database driver you need to use.

### `Driver` type [doc](https://golang.org/pkg/database/sql/driver/#Driver)

Driver is the interface that must be implemented by a database driver.

```Go
type Driver interface {
        // Open returns a new connection to the database.
        // The name is a string in a driver-specific format.
        //
        // Open may return a cached connection (one previously
        // closed), but doing so is unnecessary; the sql package
        // maintains a pool of idle connections for efficient re-use.
        //
        // The returned connection is only used by one goroutine at a
        // time.
        Open(name string) (Conn, error)
}
```

https://github.com/go-sql-driver/mysql/blob/7ac0064e822156a17a6b598957ddf5e0287f8288/driver.go

```Go
func (d MySQLDriver) Open(dsn string) (driver.Conn, error) {
    ...
}
```

https://github.com/lib/pq/blob/90697d60dd844d5ef6ff15135d0203f65d2f53b8/conn.go

```Go
// Driver is the Postgres database driver.
type Driver struct{}

// Open opens a new connection to the database. name is a connection string.
// Most users should only use it through database/sql package from the standard
// library.
func (d *Driver) Open(name string) (driver.Conn, error) {
	return Open(name)
}
```

**pq** Example

```Go
import (
	"database/sql"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	age := 21
	rows, err := db.Query("SELECT name FROM users WHERE age = $1", age)
	…
}
```

create a sql.DB you can use sql.Open(). Open returns a \*sql.DB.

`sql.Open()` does not establish any connection to the database and does not validate the driver connection parameters. Instead, it just prepares the database abstraction for later use. The first real connection to the underlying data store will be lazy to build the first time it is needed. If you want to check if the database is available immediately (for example, check if you can establish a network connection and log in), use `db.Ping()`.
