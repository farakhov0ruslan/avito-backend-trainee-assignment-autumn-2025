-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pr_reviewers (
    pr_id VARCHAR(255) NOT NULL,
    reviewer_id VARCHAR(255) NOT NULL,
    assigned_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (pr_id, reviewer_id),
    CONSTRAINT fk_pr_reviewers_pr FOREIGN KEY (pr_id) REFERENCES pull_requests(id) ON DELETE CASCADE,
    CONSTRAINT fk_pr_reviewers_user FOREIGN KEY (reviewer_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_pr_reviewers_reviewer_id ON pr_reviewers(reviewer_id);
CREATE INDEX idx_pr_reviewers_assigned_at ON pr_reviewers(assigned_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pr_reviewers CASCADE;
-- +goose StatementEnd
