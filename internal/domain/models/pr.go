package models

import (
	"fmt"
	"time"
)

// PRStatus представляет статус Pull Request
type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

// IsValid проверяет корректность статуса
func (s PRStatus) IsValid() bool {
	return s == PRStatusOpen || s == PRStatusMerged
}

type PullRequest struct {
	ID                string     `json:"pull_request_id" db:"id"`
	Name              string     `json:"pull_request_name" db:"name"`
	AuthorID          string     `json:"author_id" db:"author_id"`
	Status            PRStatus   `json:"status" db:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         time.Time  `json:"createdAt,omitempty" db:"created_at"`
	MergedAt          *time.Time `json:"mergedAt,omitempty" db:"merged_at"`
}

// IsMerged проверяет, является ли PR merged
func (pr *PullRequest) IsMerged() bool {
	return pr.Status == PRStatusMerged
}

// IsReviewerAssigned проверяет, назначен ли пользователь ревьювером
func (pr *PullRequest) IsReviewerAssigned(userID string) bool {
	for _, reviewerID := range pr.AssignedReviewers {
		if reviewerID == userID {
			return true
		}
	}
	return false
}

// String возвращает строковое представление PR для логирования
func (pr *PullRequest) String() string {
	return fmt.Sprintf("PullRequest{ID: %s, Name: %s, AuthorID: %s, Status: %s, Reviewers: %v}",
		pr.ID, pr.Name, pr.AuthorID, pr.Status, pr.AssignedReviewers)
}
