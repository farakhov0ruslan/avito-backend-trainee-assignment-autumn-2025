package response

import "time"

// PullRequestResponse ?>;=>5 ?@54AB02;5=85 PR
type PullRequestResponse struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

// PullRequestShortResponse :@0B:>5 ?@54AB02;5=85 PR (4;O A?8A:>2)
type PullRequestShortResponse struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

// CreatePRResponse >15@B:0 4;O POST /pullRequest/create
type CreatePRResponse struct {
	PR PullRequestResponse `json:"pr"`
}

// MergePRResponse >15@B:0 4;O POST /pullRequest/merge
type MergePRResponse struct {
	PR PullRequestResponse `json:"pr"`
}

// ReassignReviewerResponse >15@B:0 4;O POST /pullRequest/reassign
type ReassignReviewerResponse struct {
	PR         PullRequestResponse `json:"pr"`
	ReplacedBy string              `json:"replaced_by"`
}
