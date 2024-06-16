package easysqlite

import (
	"context"
	"database/sql"

	"github.com/pav5000/easy-sqlite/internal/errors"
	"github.com/pav5000/easy-sqlite/internal/txctx"
)

// DoInTx starts transaction and guarantees that all queries in easysqlite methods
// will use it as long as you use the context provided into the callback function.
func (s *EasySqlite) DoInTx(ctx context.Context, callback func(ctx context.Context) error) error {
	transaction, err := s.db.BeginTxx(ctx, &sql.TxOptions{
		ReadOnly:  false,
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return errors.Wrp(err, "db.BeginTxx")
	}

	defer func() {
		_ = transaction.Rollback()
	}()

	// injecting transaction into context so all query methods will use it
	ctx = txctx.Inject(ctx, transaction)

	err = callback(ctx)
	if err != nil {
		return err
	}

	// Committing only if there is no error
	return errors.Wrp(transaction.Commit(), "commit")
}
