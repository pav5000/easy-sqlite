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

You need to store mirations inside *.sql files in a folder relative to the source file which calls `easysqlite.New`. Each SQL file must have a number prefix followed by an underscore. Because migrations will be applied in the order of these number prefixes.

```
$ tree migrations/
migrations/
├── 001_initial.sql
└── 002_added_age.sql

1 directory, 2 files
```

You may check out the example main.go and migrations in this repository. The format of the migrations should be according to goose guidelines because goose library is used under the hood.

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

The `Down` clause is needed if you want rollback your migrations in the future. My library doesn't support rollbacks yet but you can do it with goose CLI-tool.

## Should I keep migration sql files near the app binary?

No, after compiling migrations get embedded into the executable binary (with the help of go's [embed](https://pkg.go.dev/embed) package). So you don't need to include migration files into the deploy container along with your app. They're needed only for build.
