package db

import (
	"context"
	"database/sql"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createFakeTransfers(t *testing.T) Transfer {
	account, er := testConnection.GetLastAccount(context.Background())
	require.NoError(t, er)
	require.NotEmpty(t, account)

	id, _ := faker.RandomInt(1, int(account.ID), 2)
	amunt, _ := faker.RandomInt(1, 10000, 1)

	arg := CreateTransferParams{Amount: int64(amunt[0]),
		FromAccountID: sql.NullInt64{Int64: int64(id[0]), Valid: true},
		ToAccountID:   sql.NullInt64{Int64: int64(id[1]), Valid: true},
	}
	sr, er := testConnection.CreateTransfer(context.Background(), arg)

	require.NoError(t, er)
	require.NotEmpty(t, sr)
	trans, er := testConnection.GetLastTransfer(context.Background())
	require.NoError(t, er)
	require.NotEmpty(t, trans)

	require.Equal(t, trans.Amount, arg.Amount)
	require.Equal(t, trans.FromAccountID, arg.FromAccountID)
	require.Equal(t, trans.ToAccountID, arg.ToAccountID)

	require.NotZero(t, trans.ID)
	require.NotZero(t, trans.CreatedAt)

	return trans
}
func TestQueries_CreateTransfers(t *testing.T) {
	createFakeTransfers(t)
}
func TestQueries_GetTransfersById(t *testing.T) {
	trans1 := createFakeTransfers(t)
	trans2, er := testConnection.GetTransferById(context.Background(), trans1.ID)
	require.NoError(t, er)
	require.NotEmpty(t, trans2)
	require.Equal(t, trans1.ID, trans2.ID)
	require.Equal(t, trans1.ToAccountID, trans2.ToAccountID)
	require.Equal(t, trans1.FromAccountID, trans2.FromAccountID)

	require.WithinDuration(t, trans2.CreatedAt, trans1.CreatedAt, time.Second)

}
func TestQueries_UpdateTransfers(t *testing.T) {
	amount, _ := faker.RandomInt(100, 1000000, 1)
	trans := createFakeTransfers(t)
	arg := UpdateTransferParams{ID: trans.ID, Amount: int64(amount[0])}
	er := testConnection.UpdateTransfer(context.Background(), arg)
	require.NoError(t, er)

}

func TestQueries_DeleteTransfers(t *testing.T) {
	trans := createFakeTransfers(t)
	er := testConnection.DeleteTransfer(context.Background(), trans.ID)
	require.NoError(t, er)
}
func TestQueries_ListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createFakeTransfers(t)
	}
	arg := ListTransfersParams{Limit: 5, Offset: 5}
	transfersList, er := testConnection.ListTransfers(context.Background(), arg)
	require.NoError(t, er)
	require.Len(t, transfersList, 5)
	for _, ac := range transfersList {
		require.NotEmpty(t, ac)
	}
}
func TestQueries_GetLastTransfer(t *testing.T) {
	trans1 := createFakeTransfers(t)
	trans2, er := testConnection.GetLastTransfer(context.Background())
	require.NoError(t, er)
	require.NotEmpty(t, trans2)
	require.Equal(t, trans1.ID, trans2.ID)
	require.Equal(t, trans1.ToAccountID, trans2.ToAccountID)
	require.Equal(t, trans1.FromAccountID, trans2.FromAccountID)

	require.WithinDuration(t, trans2.CreatedAt, trans1.CreatedAt, time.Second)
}
