package main

import (
	"context"
	"embed"

	easysqlite "github.com/pav5000/easy-sqlite"
)

// embeding migrations folder into executable binary
// the path should be relative to your source file
//
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
