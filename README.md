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

You may check out the example [main.go](https://github.com/pav5000/easy-sqlite/blob/master/cmd/example/main.go) and [migrations](https://github.com/pav5000/easy-sqlite/tree/master/cmd/example/migrations) in this repository. The format of the migrations should be according to goose guidelines because goose library is used under the hood.

https://github.com/pressly/goose

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

The `Down` clause is needed if you want to rollback your migrations in the future. My library doesn't support rollbacks yet but you can do it with goose CLI-tool.

### Should I keep migration sql files near the app binary?

No, after compiling migrations get embedded into the executable binary (with the help of go's [embed](https://pkg.go.dev/embed) package). So you don't need to include migration files into the deploy container along with your app. They're needed only for build.

### When are migrations get applied?

When you call `easysqlite.New(...)` goose checks if there are some migrations that weren't applied yet. If there are, goose applies them. If there is no database file, it will be created and all migrations applied to it one-by-one.

I think it's a good way for small apps and pet projects to apply migrations on start.
