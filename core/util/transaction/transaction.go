package transaction

import (
	"database/sql"
	"errors"
	"licor_model/core/server/shared"
)

func RunInTx(fn func(tx *sql.Tx) error) error {
	tx, err := shared.GetDB().Begin()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err == nil {
		return tx.Commit()
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}

	return err
}
