# Easy SQLite

## This is the library for you if

- You want to reduce boilerplate of SQLite init in your pet project
- You want automatic migrations for your app

## How to use

```go
//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
    // Creating db connection
    db, err := easysqlite.New("db.sqlite", embedMigrations, "migrations")
    if err != nil {
    	panic(err)
    }

    // Inserting records
    ctx := context.Background()
    _, err = db.ExecContext(ctx, `INSERT INTO users (name,age) VALUES(?,?)`, "John", 23)
    if err != nil {
    	panic(err)
    }

    // User represents one row of the table "users"
    type User struct {
    	ID   int64  `db:"id"`
    	Name string `db:"name"`
    	Age  int    `db:"age"`
    }

    // Getting one row
    var user User
    err = db.GetContext(ctx, &user, `SELECT id,name,age FROM users WHERE id=?`, 1)
    if err != nil {
    	panic(err)
    }

    // Selecting many rows
    var users []User
    err = db.SelectContext(ctx, &users, `SELECT id,name,age FROM users`)
    if err != nil {
    	panic(err)
    }
}
```

You need to store mirations inside *.sql files in a folder relative to the source file which contains `var embedMigrations embed.FS`. Each SQL file must have a number prefix followed by an underscore. Because migrations will be applied in the order of these number prefixes.

```
$ tree migrations/
migrations/
├── 001_initial.sql
└── 002_added_age.sql

1 directory, 2 files
```

The formatting convention used for the migrations are taken from the [goose](https://github.com/pressly/goose) library, since it's being used under the hood.
Check out the [main.go](https://github.com/pav5000/easy-sqlite/blob/master/cmd/example/main.go) file and the [migrations](https://github.com/pav5000/easy-sqlite/tree/master/cmd/example/migrations) folder which are within this repository as a reference example

The contents of a simple migration file:

```sql
-- +goose Up
CREATE TABLE users (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name STRING NOT NULL
);

-- +goose Down
DROP TABLE users;
```

The Down clause is needed if you want to rollback your migrations in the future. Note: Rollbacks aren't currently supported, for the meantime if needed you can use the goose CLI-tool.

### Should I keep migration sql files near the app binary?

No, after compiling migrations get embedded into the executable binary (with the help of go's [embed](https://pkg.go.dev/embed) package). So you don't need to include migration files into the deploy container along with your app. They're needed only for the build.

### When do migrations get applied?

When you call `easysqlite.New(...)` goose checks if there are any migrations that weren't applied yet. If there are, goose applies them. If there is no database file, it will be created and all migrations will be applied to it one-by-one.

Personally I think this a good way for small apps and pet projects to apply migrations on startup.

## Transactions

```go
err = db.DoInTx(ctx, func(ctx context.Context) error {
    transferAmount := 200
    userFrom := 100
    userTo := 101

    var currentBalance int
    err := db.GetContext(ctx, &currentBalance,
    	`SELECT balance FROM users WHERE id=?`,
    	userFrom)
    if err != nil {
    	return err
    }

    if currentBalance < transferAmount {
    	return errors.New("insufficient funds")
    }

    _, err = db.ExecContext(ctx,
    	`UPDATE users SET balance=balance-? WHERE id=?`,
    	transferAmount, userFrom)
    if err != nil {
    	return err
    }

    _, err = db.ExecContext(ctx,
    	`UPDATE users SET balance=balance+? WHERE id=?`,
    	transferAmount, userTo)
    if err != nil {
    	return err
    }

    return nil
})
```

`DoInTx` initiates a transaction with LevelSerializable and commits if the provided callback function returns nil. If it returns any error, the transaction is rolled back.

It is unnecessary to pass the tx object to query methods; they will retrieve the tx object from the context. Ensure that the context obtained in the callback function is passed to all query methods and that only the query methods exported by easysqlite are utilized.

Utilizing transactions in this manner is advantageous because, for instance, when employing [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html), it is challenging to maintain the domain layer free from database-specific implementations when transactions are required. When the transaction is automatically obtained from the context, it liberates the domain layer from implementation-specific imports. You can simply encapsulate DoInTx using an interface and invoke your repository methods.

Even if clean architecture is not employed, this approach remains beneficial as it eliminates the need to manually manage transaction beginnings and rollbacks, ensuring that you never forget to pass a tx object.
