package request

// TeamMemberRequest ?@54AB02;O5B CG0AB=8:0 :><0=4K 2 70?@>A5
type TeamMemberRequest struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// CreateTeamRequest 4;O POST /team/add
// 0;840F8O:
// - Handler: ?@>25@:0 GB> ?>;O =5 ?CABK5, :>@@5:B=K9 JSON
// - Service: team_name C=8:0;5=, members =5 ?CAB>9 <0AA82
type CreateTeamRequest struct {
	TeamName string              `json:"team_name"`
	Members  []TeamMemberRequest `json:"members"`
}

// GetTeamRequest 4;O GET /team/get?team_name=...
// 0;840F8O:
// - Handler: ?@>25@:0 GB> team_name =5 ?CAB>9
// - Service: ?@>25@:0 ACI5AB2>20=8O :><0=4K
type GetTeamRequest struct {
	TeamName string `json:"team_name"`
}
