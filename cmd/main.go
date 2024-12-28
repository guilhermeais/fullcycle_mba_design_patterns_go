package main

import (
	"context"
	"log"

	usecase "invoices/internal/app/application/usecases"
	httpHandlers "invoices/internal/app/infrastructure/http"
	repository "invoices/internal/app/infrastructure/repository"
	"net/http"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	pgConnection, err := repository.MakePGConnection()
	if err != nil {
		log.Fatalf("error on creating the pg connection: %v", err)
	}
	defer pgConnection.Close(context.Background())
	contractRepository := repository.NewPSQLContractRepository(*pgConnection)
	generateInvoices := usecase.NewGenerateInvoices(contractRepository)
	generateInvoicesHandler := &httpHandlers.GenerateInvoicesHandler{UseCase: generateInvoices}

	http.Handle("/generate-invoices", generateInvoicesHandler)
	log.Println("Servidor HTTP iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
