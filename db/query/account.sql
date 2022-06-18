-- name: CreateAccount :execresult
 INSERT INTO accounts (
        owner,
        balance,
        currency
 )VALUES (
        ?, ?, ?
) ON DUPLICATE KEY UPDATE balance=balance;
-- name: GetLastAccount :one
SELECT * FROM accounts ORDER BY id DESC LIMIT 1;

-- name: GetAccountById :one
SELECT * FROM accounts WHERE id=? LIMIT 1;
-- name: GetAccountByIdForUpdate :one
SELECT * FROM accounts
WHERE id=? LIMIT 1 FOR UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts ORDER BY id LIMIT ? OFFSET ?;
-- name: UpdateAccount :exec
UPDATE accounts
SET balance=balance+?
WHERE id=?;
-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id=?;