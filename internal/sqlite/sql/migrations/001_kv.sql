-- +goose Up
CREATE TABLE kv 
(
    key text PRIMARY KEY, 
    val text NOT NULL
);

-- +goose Down
DROP TABLE kv;