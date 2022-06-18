package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_TransferTx(t *testing.T) {
	store := NewStore(dbTest)
	account1 := createFakeAccount(t)
	account2 := createFakeAccount(t)
	amont := int64(10)
	//fmt.Println(account1.Balance, account2.Balance)
	//var errChannel []error
	//var resultChannel []TransferTxResult
	errChannel := make(chan error)
	resultChannel := make(chan TransferTxResult)
	count := 5
	for i := 0; i < count; i++ {
		txName := fmt.Sprintf("tx %d", i)
		go func() {
			ctx := context.WithValue(context.Background(), TxKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				Amount:        amont,
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
			})
			//errChannel = append(errChannel, err)
			//resultChannel = append(resultChannel, result)
			errChannel <- err
			resultChannel <- result
		}()
	}
	check := make(map[int]bool)
	for i := 0; i < count; i++ {
		err := <-errChannel
		result := <-resultChannel
		require.NoError(t, err)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID.Int64, account1.ID)
		require.Equal(t, transfer.ToAccountID.Int64, account2.ID)
		require.Equal(t, amont, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		fromentry := result.FromEntry
		require.NotEmpty(t, fromentry)
		require.Equal(t, fromentry.AccountID.Int64, account1.ID)
		require.Equal(t, amont, -1*fromentry.Amount)
		require.NotZero(t, fromentry.ID)
		require.NotZero(t, fromentry.CreatedAt)

		toentry := result.ToEntry
		require.NotEmpty(t, toentry)
		require.Equal(t, toentry.AccountID.Int64, account2.ID)
		require.Equal(t, amont, toentry.Amount)
		require.NotZero(t, toentry.ID)
		require.NotZero(t, toentry.CreatedAt)

		toaccount := result.ToAccount
		require.NotEmpty(t, toaccount)
		require.Equal(t, toaccount.ID, account2.ID)
		fromaccount := result.FromAccount
		diff := toaccount.Balance - account2.Balance
		diff2 := account1.Balance - fromaccount.Balance
		require.Equal(t, diff2, diff)
		require.True(t, diff > 0)

		require.True(t, diff%amont == 0)
		k := int(diff / amont)
		require.True(t, k >= 1 && k <= count)
		require.NotContains(t, check, k)
		check[k] = true

	}
}

func TestStore_TransferTxDeadLock(t *testing.T) {
	store := NewStore(dbTest)
	account1 := createFakeAccount(t)
	account2 := createFakeAccount(t)
	amont := int64(10)
	//fmt.Println(account1.Balance, account2.Balance)
	//var errChannel []error
	//var resultChannel []TransferTxResult
	errChannel := make(chan error)

	count := 5
	for i := 0; i < count; i++ {
		txName := fmt.Sprintf("tx %d", i)
		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			ctx := context.WithValue(context.Background(), TxKey, txName)
			_, err := store.TransferTx(ctx, TransferTxParams{
				Amount:        amont,
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
			})
			//errChannel = append(errChannel, err)
			//resultChannel = append(resultChannel, result)
			errChannel <- err

		}()
	}

	for i := 0; i < count; i++ {
		err := <-errChannel

		require.NoError(t, err)

	}
}
