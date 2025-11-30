package service

import (
	"context"
	"fmt"
	"time"

	"github.com/andreychano/api-service/internal/domain"
)

// AnswerService handles business logic related to answers.
type AnswerService struct {
	answerRepo   domain.AnswerRepository
	questionRepo domain.QuestionRepository
}

func NewAnswerService(answerRepo domain.AnswerRepository, questionRepo domain.QuestionRepository) *AnswerService {
	return &AnswerService{
		answerRepo:   answerRepo,
		questionRepo: questionRepo,
	}
}

// CreateAnswer validates existence of the question and stores the answer.
func (s *AnswerService) CreateAnswer(ctx context.Context, questionID int, userID, text string) (*domain.Answer, error) {
	if text == "" {
		return nil, fmt.Errorf("answer text cannot be empty")
	}
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Business Rule: Ensure the question exists before answering.
	if _, err := s.questionRepo.GetByID(ctx, questionID); err != nil {
		return nil, fmt.Errorf("question not found or error: %w", err)
	}

	a := &domain.Answer{
		QuestionID: questionID,
		UserID:     userID,
		Text:       text,
		CreatedAt:  time.Now(),
	}

	if err := s.answerRepo.Create(ctx, a); err != nil {
		return nil, fmt.Errorf("failed to create answer: %w", err)
	}

	return a, nil
}

// DeleteAnswer removes an answer by ID.
func (s *AnswerService) DeleteAnswer(ctx context.Context, id int) error {
	return s.answerRepo.Delete(ctx, id)
}
