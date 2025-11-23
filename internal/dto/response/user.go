package response

// UserResponse ?@54AB02;O5B ?>;L7>20B5;O 2 >B25B5
type UserResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

// SetUserActiveResponse >15@B:0 4;O POST /users/setIsActive
type SetUserActiveResponse struct {
	User UserResponse `json:"user"`
}

// GetUserReviewsResponse 4;O GET /users/getReview
type GetUserReviewsResponse struct {
	UserID       string                     `json:"user_id"`
	PullRequests []PullRequestShortResponse `json:"pull_requests"`
}
