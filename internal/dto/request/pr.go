package request

// CreatePRRequest  POST /pullRequest/create
type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

// MergePRRequest  POST /pullRequest/merge
type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

// ReassignReviewerRequest  POST /pullRequest/reassign
type ReassignReviewerRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}
