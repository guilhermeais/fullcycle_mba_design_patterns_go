package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	usecase "invoices/internal/app/application/usecases"
	domain "invoices/internal/app/domain/entities"
	httpHandlers "invoices/internal/app/infrastructure/http"
	"invoices/internal/app/infrastructure/repository"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGenerateInvoicesHandler(t *testing.T) {
	err := godotenv.Load("../../../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	pgConnection, err := repository.MakePGConnection()
	if err != nil {
		log.Fatalf("error on creating the pg connection: %v", err)
	}
	defer pgConnection.Close(context.Background())
	contractRepository := repository.NewPSQLContractRepository(*pgConnection)
	generateInvoices := usecase.NewGenerateInvoices(contractRepository)
	generateInvoicesHandler := &httpHandlers.GenerateInvoicesHandler{UseCase: generateInvoices}
	server := httptest.NewServer((generateInvoicesHandler))
	defer server.Close()

	input := usecase.GenerateInvoicesInput{
		Year:  2024,
		Month: 12,
		Type:  domain.InvoiceTypeAccrual,
	}
	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal input: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, server.URL+"/generate-invoices", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	output := []usecase.GenerateInvoicesOutput{}
	err = json.NewDecoder(resp.Body).Decode(&output)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	assert.Len(t, output, 1)
	assert.Equal(t, usecase.GenerateInvoicesOutput{Date: "2024-12-19", Amount: 500}, output[0])
}
