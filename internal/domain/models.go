package domain

import (
	"time"
)

// Question represents a user-submitted question.
type Question struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	Answers   []Answer  `json:"answers,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
}

// Answer represents a response to a specific question.
type Answer struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	QuestionID int       `json:"question_id"`
	UserID     string    `json:"user_id"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
}
