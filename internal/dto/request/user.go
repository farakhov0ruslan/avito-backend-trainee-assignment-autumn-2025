package request

// SetUserActiveRequest 4;O POST /users/setIsActive
// 0;840F8O:
// - Handler: ?@>25@:0 GB> ?>;O =5 ?CABK5, is_active MB> bool
// - Service: ?@>25@:0 GB> ?>;L7>20B5;L ACI5AB2C5B
type SetUserActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

// GetUserReviewsRequest 4;O GET /users/getReview?user_id=...
// 0;840F8O:
// - Handler: ?@>25@:0 GB> user_id =5 ?CAB>9
// - Service: ?@>25@:0 GB> ?>;L7>20B5;L ACI5AB2C5B (>?F8>=0;L=>, <>65B 25@=CBL ?CAB>9 <0AA82)
type GetUserReviewsRequest struct {
	UserID string `json:"user_id"`
}
