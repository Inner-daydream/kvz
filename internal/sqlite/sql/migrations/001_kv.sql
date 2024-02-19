-- +goose Up
CREATE TABLE kv
(
    "key" TEXT PRIMARY KEY,
    val TEXT NOT NULL
);

CREATE TABLE hooks
(
    name TEXT PRIMARY KEY,
    script TEXT,
    is_file BOOLEAN DEFAULT FALSE NOT NULL,
    filepath TEXT
);

CREATE TABLE key_hooks
(
    "key" TEXT NOT NULL,
    hook TEXT NOT NULL,
    FOREIGN KEY ("key") REFERENCES kv ("key"),
    FOREIGN KEY (hook) REFERENCES hooks ("name")
);

-- +goose Down
DROP TABLE kv;
DROP TABLE hooks;
DROP TABLE key_hooks;
