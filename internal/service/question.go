package service

import (
	"context"
	"fmt"
	"time"

	"github.com/andreychano/api-service/internal/domain"
)

// QuestionService handles business logic related to questions.
type QuestionService struct {
	repo domain.QuestionRepository
}

func NewQuestionService(repo domain.QuestionRepository) *QuestionService {
	return &QuestionService{repo: repo}
}

// CreateQuestion validates input and stores a new question.
func (s *QuestionService) CreateQuestion(ctx context.Context, text string) (*domain.Question, error) {
	if text == "" {
		return nil, fmt.Errorf("question text cannot be empty")
	}
	if len(text) > 255 {
		return nil, fmt.Errorf("question text is too long")
	}

	q := &domain.Question{
		Text:      text,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, q); err != nil {
		return nil, fmt.Errorf("failed to create question: %w", err)
	}

	return q, nil
}

// GetAllQuestions retrieves all questions from storage.
func (s *QuestionService) GetAllQuestions(ctx context.Context) ([]domain.Question, error) {
	return s.repo.GetAll(ctx)
}

// GetQuestionWithAnswers retrieves a specific question and its associated answers.
func (s *QuestionService) GetQuestionWithAnswers(ctx context.Context, id int) (*domain.Question, error) {
	return s.repo.GetByID(ctx, id)
}

// DeleteQuestion removes a question and cascades delete to its answers.
func (s *QuestionService) DeleteQuestion(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
