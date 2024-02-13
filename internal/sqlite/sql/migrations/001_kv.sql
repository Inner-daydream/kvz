-- +goose Up
CREATE TABLE kv
(
    "key" text PRIMARY KEY,
    val text NOT NULL
);

CREATE TABLE hooks
(
    name text PRIMARY KEY,
    script text NOT NULL
);

CREATE TABLE key_hooks
(
    "key" text NOT NULL,
    hook text NOT NULL,
    FOREIGN KEY ("key") REFERENCES kv ("key"),
    FOREIGN KEY (hook) REFERENCES hooks ("name")
);

-- +goose Down
DROP TABLE kv;
DROP TABLE hooks;
DROP TABLE key_hooks;
