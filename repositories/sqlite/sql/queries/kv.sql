-- name: GetVal :one
SELECT
    val
FROM
    kv
WHERE
    key = ?;

-- name: SetVal :exec
INSERT OR REPLACE INTO kv (key, val)
VALUES
    (?, ?);