package models

import "fmt"

type Team struct {
	Name    string `json:"team_name" db:"name"`
	Members []User `json:"members"`
}

// GetActiveMembers возвращает только активных участников команды
func (t *Team) GetActiveMembers() []User {
	activeMembers := make([]User, 0)
	for _, member := range t.Members {
		if member.IsActive {
			activeMembers = append(activeMembers, member)
		}
	}
	return activeMembers
}

// GetActiveMembersExcept возвращает активных участников команды, исключая указанного пользователя
func (t *Team) GetActiveMembersExcept(userID string) []User {
	activeMembers := make([]User, 0)
	for _, member := range t.Members {
		if member.IsActive && member.ID != userID {
			activeMembers = append(activeMembers, member)
		}
	}
	return activeMembers
}

// String возвращает строковое представление команды для логирования
func (t *Team) String() string {
	return fmt.Sprintf("Team{Name: %s, MembersCount: %d}", t.Name, len(t.Members))
}
