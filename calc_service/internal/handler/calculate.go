package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/Alehamrom/calc_service/internal/evaluator"
	"github.com/Alehamrom/calc_service/pkg/errors"
)

type CalculateRequest struct {
	Expression string `json:"expression"`
}

type CalculateResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CalculateRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, "Invalid request payload", errors.UnprocessableEntity)
		return
	}

	if !isValidExpression(req.Expression) {
		respondWithError(w, "Expression is not valid", errors.UnprocessableEntity)
		return
	}

	result, err := evaluator.Evaluate(req.Expression)
	if err != nil {
		if err.Error() == "invalid character in expression" || err.Error() == "mismatched parentheses" || err.Error() == "unknown token" {
			respondWithError(w, "Expression is not valid", errors.UnprocessableEntity)
			return
		}
		respondWithError(w, "Internal server error", errors.InternalServerError)
		return
	}

	respondWithResult(w, fmt.Sprintf("%v", result))
}

func isValidExpression(expr string) bool {
	re := regexp.MustCompile(`^[0-9+\-*/^().\s]+$`)
	return re.MatchString(expr)
}

func respondWithError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(CalculateResponse{Error: message})
}

func respondWithResult(w http.ResponseWriter, result string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CalculateResponse{Result: result})
}
