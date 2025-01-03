package http_handlers_test

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
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}
func TestGenerateInvoicesHandler(t *testing.T) {
	t.Run("Deve gerar faturas por regime de competência via API", func(t *testing.T) {
		t.Parallel()
		server, contractFactory := makeSut(t)
		defer server.Close()

		mockedContract := domain.Contract{
			Id:          "fac05a57-7d61-4283-ab32-7696902eac44",
			Description: "prestação de serviços escolares",
			Amount:      6000,
			Periods:     12,
			Date:        time.Date(2024, time.Month(12), 18, 10, 0, 0, 0, time.UTC),
		}
		err := contractFactory.CreateContract(mockedContract)
		if err != nil {
			t.Fatalf("error creating mocked contract: %v\n", err)
		}

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
		assert.Equal(t, usecase.GenerateInvoicesOutput{Date: "2024-12-18", Amount: 500}, output[0])
	})

	t.Run("Deve gerar faturas pro regime de caixa via API", func(t *testing.T) {
		t.Parallel()
		server, contractFactory := makeSut(t)
		defer server.Close()

		mockedContract := domain.Contract{
			Id:          "fac05a57-7d61-4283-ab32-7696902eac44",
			Description: "prestação de serviços escolares",
			Amount:      6000,
			Periods:     12,
			Date:        time.Date(2024, time.Month(12), 18, 10, 0, 0, 0, time.UTC),
		}
		mockedContract.AddPayment(domain.Payment{
			Id:     "6355b223-fce0-4f7c-998a-1f027281e308",
			Amount: 6000,
			Date:   time.Date(2024, time.Month(12), 18, 10, 0, 0, 0, time.UTC),
		})
		contractFactory.CreateContract(mockedContract)
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

func makeSut(t *testing.T) (*httptest.Server, testutils.ContractFactory) {
	schema, err := testutils.CreateNewSchema()
	if err != nil {
		t.Fatalf("error on creating schema for this test: %v", err)
	}

	pgConnection, err := repository.MakePGConnectionWithUri(schema.URI)
	if err != nil {
		t.Fatalf("error on creating the pg connection: %v", err)
	}
	t.Cleanup(func() {
		pgConnection.Close(context.Background())
	})
	dbMigrator := testutils.PostgresDbMigrator{Conn: *pgConnection, Schema: schema.Schema}
	dbMigrator.MigrateDb()
	t.Cleanup(func() {
		dbMigrator.DropDb()
	})
	contractRepository := repository.NewPSQLContractRepository(*pgConnection)
	generateInvoices := usecase.NewGenerateInvoices(contractRepository, usecase.NewObserver[usecase.InvoiceGeneratedEventData]())
	generateInvoicesHandler := httpHandlers.LoggerDecorator{&httpHandlers.GenerateInvoicesHandler{UseCase: generateInvoices}}

	contractFactory := testutils.ContractFactory{Conn: pgConnection}
	return httptest.NewServer(generateInvoicesHandler), contractFactory
}

func setup() {
	err := godotenv.Load("../../../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
