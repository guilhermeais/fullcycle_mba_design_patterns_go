package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	usecase "invoices/internal/app/application/usecases"
	domain "invoices/internal/app/domain/entities"
	httpHandlers "invoices/internal/app/infrastructure/http"
	"invoices/internal/app/infrastructure/repository"
	testutils "invoices/internal/testutils"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGenerateInvoicesHandler(t *testing.T) {
	t.Run("Deve gerar faturas por regime de competência via API", func(t *testing.T) {
		server := makeSut(t)
		defer server.Close()

		body := makeBody(t, usecase.GenerateInvoicesInput{
			Year:  2024,
			Month: 12,
			Type:  domain.InvoiceTypeAccrual,
		})

		resp := postGenerateInvoices(t, server, body)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		output := decodeInvoicesResponse(t, resp)
		assert.Len(t, output, 1)
		assert.Equal(t, usecase.GenerateInvoicesOutput{Date: "2024-12-19", Amount: 500}, output[0])
	})

	t.Run("Deve gerar faturas pro regime de caixa via API", func(t *testing.T) {
		server := makeSut(t)
		defer server.Close()
		body := makeBody(t, usecase.GenerateInvoicesInput{
			Year:  2024,
			Month: 12,
			Type:  domain.InvoiceTypeCash,
		})
		resp := postGenerateInvoices(t, server, body)
		defer resp.Body.Close()
		assert.Equal(t, resp.StatusCode, http.StatusOK)
		output := decodeInvoicesResponse(t, resp)
		assert.Len(t, output, 1)
		assert.Equal(t, usecase.GenerateInvoicesOutput{Date: "2024-12-18", Amount: 6000}, output[0])
	})
}

func decodeInvoicesResponse(t *testing.T, resp *http.Response) []usecase.GenerateInvoicesOutput {
	output := []usecase.GenerateInvoicesOutput{}
	err := json.NewDecoder(resp.Body).Decode(&output)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return output
}

func postGenerateInvoices(t *testing.T, server *httptest.Server, body []byte) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, server.URL+"/generate-invoices", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	return resp
}

func makeBody(t *testing.T, input usecase.GenerateInvoicesInput) []byte {
	t.Helper()
	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal input: %v", err)
	}
	return body
}

func makeSut(t *testing.T) *httptest.Server {
	t.Helper()
	err := godotenv.Load("../../../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	pgConnection, err := repository.MakePGConnectionWithUri(testutils.PgContainer.URI)
	t.Cleanup(func() {
		pgConnection.Close(context.Background())
	})
	testutils.MigrateDb()
	if err != nil {
		log.Fatalf("error on creating the pg connection: %v", err)
	}
	contractRepository := repository.NewPSQLContractRepository(*pgConnection)
	generateInvoices := usecase.NewGenerateInvoices(contractRepository)
	generateInvoicesHandler := &httpHandlers.GenerateInvoicesHandler{UseCase: generateInvoices}
	return httptest.NewServer((generateInvoicesHandler))
}
