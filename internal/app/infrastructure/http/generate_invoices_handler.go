package http_handlers

import (
	"encoding/json"
	"fmt"
	usecase "invoices/internal/app/application/usecases"
	domain "invoices/internal/app/domain/entities"
	"net/http"
	"os"
)

type GenerateInvoicesHandler struct {
	UseCase *usecase.GenerateInvoices
}

type HttpError struct {
	Message string `json:"message"`
}

func (h *GenerateInvoicesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJsonError(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var input usecase.GenerateInvoicesInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error on decoding the request body: %v", err)
		h.writeJsonError(w, fmt.Sprintf("Invalid request body: '%v'", err), http.StatusBadRequest)
		return
	}

	if input.Type != domain.InvoiceTypeCash && input.Type != domain.InvoiceTypeAccrual {
		h.writeJsonError(w, fmt.Sprintf("Invalid invoice type: '%s'", input.Type), http.StatusBadRequest)
		return
	}

	output, err := h.UseCase.Execute(input)
	if err != nil {
		h.writeJsonError(w, err.Error(), http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *GenerateInvoicesHandler) writeJsonError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(HttpError{Message: message})
}
