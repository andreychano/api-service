package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/andreychano/api-service/internal/domain"
)

// Handler holds dependencies for HTTP requests processing.
type Handler struct {
	qService domain.QuestionService
	aService domain.AnswerService
}

// NewRouter constructs a http.Handler with all application routes defined.
// This separates routing logic from the main entry point.
func NewRouter(qService domain.QuestionService, aService domain.AnswerService) http.Handler {
	mux := http.NewServeMux()
	h := &Handler{
		qService: qService,
		aService: aService,
	}

	// Register routes using Go 1.22 pattern matching

	// Questions
	mux.HandleFunc("POST /questions", h.createQuestion)
	mux.HandleFunc("GET /questions", h.getAllQuestions)
	mux.HandleFunc("GET /questions/{id}", h.getQuestionByID)
	mux.HandleFunc("DELETE /questions/{id}", h.deleteQuestion)

	// Answers
	mux.HandleFunc("POST /questions/{id}/answers", h.createAnswer)
	mux.HandleFunc("DELETE /answers/{id}", h.deleteAnswer)

	return mux
}

// --- HANDLERS IMPLEMENTATION ---

// Request/Response structs (DTOs)
type createQuestionRequest struct {
	Text string `json:"text"`
}

type createAnswerRequest struct {
	UserID string `json:"user_id"`
	Text   string `json:"text"`
}

// createQuestion handles POST /questions
func (h *Handler) createQuestion(w http.ResponseWriter, r *http.Request) {
	var req createQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	q, err := h.qService.CreateQuestion(r.Context(), req.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // Simple error handling
		return
	}

	writeJSON(w, http.StatusCreated, q)
}

// getAllQuestions handles GET /questions
func (h *Handler) getAllQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := h.qService.GetAllQuestions(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, questions)
}

// getQuestionByID handles GET /questions/{id}
func (h *Handler) getQuestionByID(w http.ResponseWriter, r *http.Request) {
	// Go 1.22: r.PathValue retrieves wildcards from URL
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	q, err := h.qService.GetQuestionWithAnswers(r.Context(), id)
	if err != nil {
		// Ideally check if error is "not found" vs "internal error"
		http.Error(w, "question not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, q)
}

// deleteQuestion handles DELETE /questions/{id}
func (h *Handler) deleteQuestion(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.qService.DeleteQuestion(r.Context(), id); err != nil {
		http.Error(w, "failed to delete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// createAnswer handles POST /questions/{id}/answers
func (h *Handler) createAnswer(w http.ResponseWriter, r *http.Request) {
	// Parse Question ID from URL
	qIDStr := r.PathValue("id")
	qID, err := strconv.Atoi(qIDStr)
	if err != nil {
		http.Error(w, "invalid question id", http.StatusBadRequest)
		return
	}

	// Parse Body
	var req createAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Call Service
	a, err := h.aService.CreateAnswer(r.Context(), qID, req.UserID, req.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, a)
}

// deleteAnswer handles DELETE /answers/{id}
func (h *Handler) deleteAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.aService.DeleteAnswer(r.Context(), id); err != nil {
		http.Error(w, "failed to delete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// writeJSON is a helper to write JSON responses
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
