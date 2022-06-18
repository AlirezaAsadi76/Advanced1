-- name: CreateTransfer :execresult
INSERT INTO transfers (
    from_account_id,
    to_account_id,
    amount
)VALUES (
            ?, ?, ?
        );
-- name: GetLastTransfer :one
SELECT * FROM transfers ORDER BY id DESC LIMIT 1;

-- name: GetTransferById :one
SELECT * FROM transfers WHERE id=? LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers ORDER BY id LIMIT ? OFFSET ?;
-- name: UpdateTransfer :exec
UPDATE transfers
SET amount=?
WHERE id=?;
-- name: DeleteTransfer :exec
DELETE FROM transfers WHERE id=?;