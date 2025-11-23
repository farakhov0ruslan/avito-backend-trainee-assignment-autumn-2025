-- +goose Up
CREATE TABLE IF NOT EXISTS teams (
    name VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_teams_created_at ON teams(created_at);

-- +goose Down
DROP TABLE IF EXISTS teams CASCADE;
