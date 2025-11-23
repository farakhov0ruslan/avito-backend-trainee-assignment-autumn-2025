package response

// TeamMemberResponse ?@54AB02;O5B CG0AB=8:0 :><0=4K 2 >B25B5
type TeamMemberResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// TeamResponse 4;O >B25B>2 A :><0=4>9 (GET /team/get, POST /team/add)
type TeamResponse struct {
	TeamName string               `json:"team_name"`
	Members  []TeamMemberResponse `json:"members"`
}

// CreateTeamResponse >15@B:0 4;O POST /team/add
type CreateTeamResponse struct {
	Team TeamResponse `json:"team"`
}
