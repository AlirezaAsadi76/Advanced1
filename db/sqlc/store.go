package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

var TxKey = struct{}{}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db,
		Queries: New(db)}
}
func (store *Store) exectx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			return fmt.Errorf("tx error:%v,rb err:%v", err, rberr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"fromAccountID"`
	BalanceFrom   int64 `json:"balanceFrom"`
	ToAccountID   int64 `json:"toAccountID"`
	BalanceTo     int64 `json:"balanceTo"`
	Amount        int64 `json:"amount"`
}
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"fromAccount"`
	ToAccount   Account  `json:"toAccount"`
	FromEntry   Entry    `json:"fromEntry"`
	ToEntry     Entry    `json:"toEntry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.exectx(ctx, func(queries *Queries) error {
		//txName := ctx.Value(TxKey)

		//fmt.Println(txName, "create transfers")
		_, err := queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: sql.NullInt64{Int64: arg.FromAccountID, Valid: true},
			ToAccountID:   sql.NullInt64{Int64: arg.ToAccountID, Valid: true},
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		//fmt.Println(txName, "get  transfers")

		result.Transfer, err = queries.GetLastTransfer(ctx)
		if err != nil {
			return err
		}
		//fmt.Println(txName, "create entry 1")

		_, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: sql.NullInt64{Int64: arg.FromAccountID, Valid: true},
			Amount:    -1 * arg.Amount,
		})
		//fmt.Println(txName, "get entry 1")

		result.FromEntry, err = queries.GetLastEntry(ctx)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		//fmt.Println(txName, "get entry 2")
		_, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: sql.NullInt64{Int64: arg.ToAccountID, Valid: true},
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		//fmt.Println(txName, "get entry 2")
		result.ToEntry, err = queries.GetLastEntry(ctx)
		if arg.FromAccountID < arg.ToAccountID {
			err = updateMoney(ctx, queries, arg.FromAccountID, -1*arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			err = updateMoney(ctx, queries, arg.ToAccountID, -1*arg.Amount, arg.FromAccountID, arg.Amount)
			if err != nil {
				return err
			}
		}
		result.FromAccount, err = queries.GetAccountById(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}
		//fmt.Println(txName, "get account 2")

		result.ToAccount, err = queries.GetAccountById(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}
		return nil
	})

	return result, err
}
func updateMoney(ctx context.Context, q *Queries, accountID1, amount1, accountID2, amount2 int64) (err error) {
	err = q.UpdateAccount(ctx, UpdateAccountParams{
		ID:      accountID1,
		Balance: amount1,
	})
	if err != nil {
		return
	}
	err = q.UpdateAccount(ctx, UpdateAccountParams{
		ID:      accountID2,
		Balance: amount2,
	})
	if err != nil {
		return
	}
	return nil
}
