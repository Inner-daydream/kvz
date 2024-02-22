-- name: getVal :one
SELECT val
FROM kv
WHERE "key" = ?;

-- name: setVal :exec
INSERT OR REPLACE INTO kv ("key", val)
VALUES (?, ?);

-- name: deleteKey :exec
DELETE FROM kv
WHERE "key" = ?;

-- name: listKeys :many
SELECT "key" FROM kv;

-- name: addScriptHook :exec
INSERT OR REPLACE INTO hooks (name, script, is_file)
VALUES (?, ?, FALSE);

-- name: addFilePathHook :exec 
INSERT OR REPLACE INTO hooks (name, filepath, is_file)
VALUES (?, ?, TRUE);

-- name: addFileHook :exec
INSERT OR REPLACE INTO hooks (name, script, is_file)
VALUES (?, ?, TRUE);

-- name: attachHook :exec
INSERT INTO key_hooks ("key", hook)
VALUES (?, ?);

-- name: deleteHook :exec
DELETE FROM hooks
where name = ?;

-- name: listHooks :many
SELECT name FROM hooks;

-- name: keyExists :one
SELECT EXISTS(
    SELECT 1 
    FROM kv 
    WHERE "key"=?
);

-- name: hookExists :one
SELECT EXISTS(
    SELECT 1
    FROM hooks
    WHERE name=?
);

-- name: getAttachedHooks :many
SELECT h.name, h.script, h.is_file, h.filepath
FROM key_hooks kh
JOIN hooks h ON kh.hook = h.name
WHERE kh.key = ?;
