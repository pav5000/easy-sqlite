//nolint:wrapcheck,mnd,err113,funlen
package main

import (
	"context"
	"embed"
	"errors"

	easysqlite "github.com/pav5000/easy-sqlite"
)

// embeding migrations folder into executable binary
// the path should be relative to your source file
//
//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	// Creating conn connection
	conn, err := easysqlite.New("db.sqlite", embedMigrations, "migrations")
	if err != nil {
		panic(err)
	}

	// Inserting records
	ctx := context.Background()
	_, err = conn.ExecContext(ctx, `INSERT INTO users (name,balance) VALUES(?,?)`, "John", 0)
	if err != nil {
		panic(err)
	}

	// User represents one row of the table "users"
	type User struct {
		ID      int64  `db:"id"`
		Name    string `db:"name"`
		Balance int    `db:"balance"`
	}

	// Getting one row
	var user User
	err = conn.GetContext(ctx, &user, `SELECT id,name,balance FROM users WHERE id=?`, 1)
	if err != nil {
		panic(err)
	}

	// Selecting many rows
	var users []User
	err = conn.SelectContext(ctx, &users, `SELECT id,name,balance FROM users`)
	if err != nil {
		panic(err)
	}

	// Transactions
	conn.MustExecContext(ctx, `DELETE FROM users WHERE id IN(?,?)`, 100, 101)
	conn.MustExecContext(ctx,
		`INSERT INTO users (id,name,balance) VALUES(?,?,?)`,
		100, "Sam", 400)
	conn.MustExecContext(ctx,
		`INSERT INTO users (id,name,balance) VALUES(?,?,?)`,
		101, "Thomas", 100)

	err = conn.DoInTx(ctx, func(ctx context.Context) error {
		transferAmount := 200
		userFrom := 100
		userTo := 101

		var currentBalance int
		err := conn.GetContext(ctx, &currentBalance, `SELECT balance FROM users WHERE id=?`, userFrom)
		if err != nil {
			return err
		}

		if currentBalance < transferAmount {
			return errors.New("insufficient funds")
		}

		_, err = conn.ExecContext(ctx,
			`UPDATE users SET balance=balance-? WHERE id=?`,
			transferAmount, userFrom)
		if err != nil {
			return err
		}

		_, err = conn.ExecContext(ctx,
			`UPDATE users SET balance=balance+? WHERE id=?`,
			transferAmount, userTo)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}
