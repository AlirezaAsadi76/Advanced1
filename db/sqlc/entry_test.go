package db

import (
	"context"
	"database/sql"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createFakeEntry(t *testing.T) Entry {
	account, er := testConnection.GetLastAccount(context.Background())
	require.NoError(t, er)
	require.NotEmpty(t, account)

	id, _ := faker.RandomInt(1, int(account.ID), 1)
	amunt, _ := faker.RandomInt(1, 10000, 1)

	arg := CreateEntryParams{
		AccountID: sql.NullInt64{Int64: int64(id[0]), Valid: true},
		Amount:    int64(amunt[0]),
	}
	sr, er := testConnection.CreateEntry(context.Background(), arg)

	require.NoError(t, er)
	require.NotEmpty(t, sr)
	entry, er := testConnection.GetLastEntry(context.Background())
	require.NoError(t, er)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.Amount, arg.Amount)
	require.Equal(t, entry.AccountID, arg.AccountID)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}
func TestQueries_CreateEntry(t *testing.T) {
	createFakeEntry(t)
}
func TestQueries_GetEntryById(t *testing.T) {
	entry1 := createFakeEntry(t)
	entry2, er := testConnection.GetEntryById(context.Background(), entry1.ID)
	require.NoError(t, er)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.AccountID, entry2.AccountID)

	require.WithinDuration(t, entry2.CreatedAt, entry1.CreatedAt, time.Second)

}
func TestQueries_UpdateEntry(t *testing.T) {
	amount, _ := faker.RandomInt(100, 1000000, 1)
	entry := createFakeEntry(t)
	arg := UpdateEntryParams{ID: entry.ID, Amount: int64(amount[0])}
	er := testConnection.UpdateEntry(context.Background(), arg)
	require.NoError(t, er)

}

func TestQueries_DeleteEntry(t *testing.T) {
	entry := createFakeEntry(t)
	er := testConnection.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, er)
}
func TestQueries_ListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createFakeEntry(t)
	}
	arg := ListEntriesParams{Limit: 5, Offset: 5}
	entriesList, er := testConnection.ListEntries(context.Background(), arg)
	require.NoError(t, er)
	require.Len(t, entriesList, 5)
	for _, ac := range entriesList {
		require.NotEmpty(t, ac)
	}
}
func TestQueries_GetLastEntry(t *testing.T) {
	entry1 := createFakeEntry(t)
	entry2, er := testConnection.GetLastEntry(context.Background())
	require.NoError(t, er)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.AccountID, entry2.AccountID)

	require.WithinDuration(t, entry2.CreatedAt, entry1.CreatedAt, time.Second)
}
