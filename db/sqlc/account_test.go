package db

import (
	"context"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createFakeAccount(t *testing.T) Account {
	balance, _ := faker.RandomInt(100, 1000000, 1)
	arg := CreateAccountParams{
		Owner:    faker.Name(),
		Balance:  int64(balance[0]),
		Currency: faker.Currency()}
	sr, er := testConnection.CreateAccount(context.Background(), arg)

	require.NoError(t, er)
	require.NotEmpty(t, sr)
	account, er := testConnection.GetLastAccount(context.Background())
	require.NoError(t, er)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}
func TestQueries_CreateAccount(t *testing.T) {
	createFakeAccount(t)
}
func TestQueries_GetAccountById(t *testing.T) {
	accont := createFakeAccount(t)
	acconttemp, er := testConnection.GetAccountById(context.Background(), accont.ID)
	require.NoError(t, er)
	require.NotEmpty(t, acconttemp)
	require.Equal(t, acconttemp.ID, accont.ID)
	require.Equal(t, acconttemp.Owner, accont.Owner)
	require.Equal(t, acconttemp.Balance, accont.Balance)
	require.Equal(t, acconttemp.Currency, accont.Currency)

	require.WithinDuration(t, accont.CreatedAt, acconttemp.CreatedAt, time.Second)

}
func TestQueries_UpdateAccount(t *testing.T) {
	balance, _ := faker.RandomInt(100, 1000000, 1)
	account := createFakeAccount(t)
	arg := UpdateAccountParams{ID: account.ID, Balance: int64(balance[0])}
	er := testConnection.UpdateAccount(context.Background(), arg)
	require.NoError(t, er)

}

func TestQueries_DeleteAccount(t *testing.T) {
	accont := createFakeAccount(t)
	er := testConnection.DeleteAccount(context.Background(), accont.ID)
	require.NoError(t, er)
}
func TestQueries_ListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createFakeAccount(t)
	}
	arg := ListAccountsParams{Limit: 5, Offset: 5}
	accountList, er := testConnection.ListAccounts(context.Background(), arg)
	require.NoError(t, er)
	require.Len(t, accountList, 5)
	for _, ac := range accountList {
		require.NotEmpty(t, ac)
	}
}
func TestQueries_GetLastAccount(t *testing.T) {
	accont := createFakeAccount(t)
	accont1, er := testConnection.GetLastAccount(context.Background())
	require.NoError(t, er)
	require.NotEmpty(t, accont)
	require.Equal(t, accont.ID, accont1.ID)
	require.Equal(t, accont.Owner, accont1.Owner)
	require.Equal(t, accont.Balance, accont1.Balance)
	require.Equal(t, accont.Currency, accont1.Currency)

	require.WithinDuration(t, accont.CreatedAt, accont1.CreatedAt, time.Second)
}
