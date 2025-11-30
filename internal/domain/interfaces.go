package domain

import (
	"context"
)

// QuestionRepository defines the interface for question data storage.
type QuestionRepository interface {
	Create(ctx context.Context, question *Question) error
	GetAll(ctx context.Context) ([]Question, error)
	GetByID(ctx context.Context, id int) (*Question, error)
	Delete(ctx context.Context, id int) error
}

// AnswerRepository defines the interface for answer data storage.
type AnswerRepository interface {
	Create(ctx context.Context, answer *Answer) error
	GetByID(ctx context.Context, id int) (*Answer, error)
	Delete(ctx context.Context, id int) error
}

// QuestionService handles business logic and data validation for questions.
type QuestionService interface {
	CreateQuestion(ctx context.Context, text string) (*Question, error)
	GetAllQuestions(ctx context.Context) ([]Question, error)
	GetQuestionWithAnswers(ctx context.Context, id int) (*Question, error)
	DeleteQuestion(ctx context.Context, id int) error
}

// AnswerService handles business logic for answers.
type AnswerService interface {
	CreateAnswer(ctx context.Context, questionID int, userID, text string) (*Answer, error)
	DeleteAnswer(ctx context.Context, id int) error
}
