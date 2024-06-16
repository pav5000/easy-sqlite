package easysqlite

import (
	"io/fs"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pav5000/easy-sqlite/internal/errors"
	"github.com/pressly/goose/v3"
)

type EasySqlite struct {
	db *sqlx.DB
}

// New opens the db file or creates it if not exists
// Migrations should be embedded by the go "embed" package
// Check out the example of migrations in the cmd/example folder
// The migration tool used here is https://github.com/pressly/goose
func New(path string, migrations fs.FS, dirName string) (*EasySqlite, error) {
	service, err := createDbService(path)
	if err != nil {
		return nil, errors.Wrp(err, "creating db service")
	}

	err = goose.SetDialect("sqlite3")
	if err != nil {
		return nil, errors.Wrp(err, "goose.SetDialect")
	}

	goose.SetBaseFS(migrations)
	err = goose.Up(service.db.DB, "migrations")
	if err != nil {
		return nil, errors.Wrp(err, "goose.Up")
	}

	return service, nil
}

func createDbService(path string) (*EasySqlite, error) {
	conn, err := sqlx.Connect("sqlite3", path)
	if err != nil {
		return nil, errors.Wrp(err, "sql.Open")
	}

	return &EasySqlite{
		db: conn,
	}, nil
}
