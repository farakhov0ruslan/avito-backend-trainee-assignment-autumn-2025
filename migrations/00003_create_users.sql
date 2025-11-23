-- +goose Up

CREATE TABLE IF NOT EXISTS users
(
    id
    VARCHAR
(
    255
) PRIMARY KEY,
    username VARCHAR
(
    255
) NOT NULL,
    team_name VARCHAR
(
    255
) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW
(
),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW
(
),
    CONSTRAINT fk_users_team FOREIGN KEY
(
    team_name
) REFERENCES teams
(
    name
) ON DELETE CASCADE
    );

CREATE INDEX idx_users_team_name ON users (team_name);
CREATE INDEX idx_users_is_active ON users (is_active);
CREATE INDEX idx_users_team_active ON users (team_name, is_active);


-- +goose Down
DROP TABLE IF EXISTS users CASCADE;
