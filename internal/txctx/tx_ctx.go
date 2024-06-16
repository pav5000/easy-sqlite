package txctx

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type key struct{}

func Inject(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, key{}, tx)
}

func Extract(ctx context.Context) *sqlx.Tx {
	tx, ok := ctx.Value(key{}).(*sqlx.Tx)
	if !ok {
		return nil
	}

	return tx
}
