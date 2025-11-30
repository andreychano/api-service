package postgres

import (
	"context"

	"github.com/andreychano/api-service/internal/domain"
	"gorm.io/gorm"
)

type QuestionRepo struct {
	db *gorm.DB
}

func NewQuestionRepo(db *gorm.DB) *QuestionRepo {
	return &QuestionRepo{db: db}
}

func (r *QuestionRepo) Create(ctx context.Context, question *domain.Question) error {
	return r.db.WithContext(ctx).Create(question).Error
}

func (r *QuestionRepo) GetAll(ctx context.Context) ([]domain.Question, error) {
	var questions []domain.Question
	err := r.db.WithContext(ctx).Find(&questions).Error
	return questions, err
}

func (r *QuestionRepo) GetByID(ctx context.Context, id int) (*domain.Question, error) {
	var question domain.Question
	err := r.db.WithContext(ctx).Preload("Answers").First(&question, id).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *QuestionRepo) Delete(ctx context.Context, id int) error {
	// Delete the question. Associated answers will be automatically deleted via database cascade.
	return r.db.WithContext(ctx).Delete(&domain.Question{}, id).Error
}

type AnswerRepo struct {
	db *gorm.DB
}

func NewAnswerRepo(db *gorm.DB) *AnswerRepo {
	return &AnswerRepo{db: db}
}

func (r *AnswerRepo) Create(ctx context.Context, answer *domain.Answer) error {
	return r.db.WithContext(ctx).Create(answer).Error
}

func (r *AnswerRepo) GetByID(ctx context.Context, id int) (*domain.Answer, error) {
	var answer domain.Answer
	err := r.db.WithContext(ctx).First(&answer, id).Error
	return &answer, err
}

func (r *AnswerRepo) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.Answer{}, id).Error
}
