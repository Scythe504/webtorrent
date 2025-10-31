-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS videos (
    id TEXT PRIMARY KEY,
    magnet_link TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('processing', 'downloading', 'downloaded', 'failed')),
    file_path TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS videos;
-- +goose StatementEnd
