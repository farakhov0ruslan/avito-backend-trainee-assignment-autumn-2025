package models

import "fmt"

type User struct {
	ID       string `json:"user_id" db:"id"`
	Username string `json:"username" db:"username"`
	TeamName string `json:"team_name" db:"team_name"`
	IsActive bool   `json:"is_active" db:"is_active"`
}

// String возвращает строковое представление пользователя для логирования
func (u *User) String() string {
	return fmt.Sprintf("User{ID: %s, Username: %s, TeamName: %s, IsActive: %t}",
		u.ID, u.Username, u.TeamName, u.IsActive)
}
