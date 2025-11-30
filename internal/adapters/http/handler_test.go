package http

import (
	"context"
	"errors"
	"testing"

	"github.com/andreychano/api-service/internal/domain"
	"github.com/andreychano/api-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- MOCK REPOSITORY ---
type mockQuestionRepo struct {
	createFunc func(ctx context.Context, q *domain.Question) error
}

func (m *mockQuestionRepo) Create(ctx context.Context, q *domain.Question) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, q)
	}
	return nil
}
func (m *mockQuestionRepo) GetAll(ctx context.Context) ([]domain.Question, error) { return nil, nil }
func (m *mockQuestionRepo) GetByID(ctx context.Context, id int) (*domain.Question, error) {
	return nil, nil
}
func (m *mockQuestionRepo) Delete(ctx context.Context, id int) error { return nil }

// --- TESTS ---

func TestQuestionService_CreateQuestion(t *testing.T) {
	// Table-Driven Tests
	tests := []struct {
		name           string
		inputText      string
		mockSetup      func() *mockQuestionRepo
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:           "Error_EmptyText",
			inputText:      "",
			mockSetup:      func() *mockQuestionRepo { return &mockQuestionRepo{} },
			expectError:    true,
			expectedErrMsg: "question text cannot be empty",
		},
		{
			name:      "Success_ValidText",
			inputText: "Is this a valid question?",
			mockSetup: func() *mockQuestionRepo {
				return &mockQuestionRepo{
					createFunc: func(ctx context.Context, q *domain.Question) error {
						q.ID = 1
						return nil
					},
				}
			},
			expectError: false,
		},
		{
			name:      "Error_RepoFailure",
			inputText: "This text is valid",
			mockSetup: func() *mockQuestionRepo {
				return &mockQuestionRepo{
					createFunc: func(ctx context.Context, q *domain.Question) error {
						return errors.New("database connection error")
					},
				}
			},
			expectError:    true,
			expectedErrMsg: "database connection error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := tt.mockSetup()
			s := service.NewQuestionService(mockRepo)

			// Act
			result, err := s.CreateQuestion(context.Background(), tt.inputText)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				require.NoError(t, err) // require stops test on failure
				assert.NotNil(t, result)
				assert.Equal(t, tt.inputText, result.Text)
				assert.NotZero(t, result.ID, "ID should be assigned")
			}
		})
	}
}
