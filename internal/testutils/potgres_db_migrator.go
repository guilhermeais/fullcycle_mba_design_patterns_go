package testutils

import (
	"context"
	repository "invoices/internal/app/infrastructure/repository"
	"log"
	"os"
)

func MigrateDb() {
	ctx := context.Background()
	pgConn, err := repository.MakePGConnectionWithUri(PgContainer.URI)
	if err != nil {
		log.Fatalf("error migrating db: %v", err)
	}
	defer pgConn.Close(ctx)
	sqlMigrationFilePath := os.Getenv("TEST_SQL_MIGRATION_PATH")

	log.Printf("migrating with file '%s'", sqlMigrationFilePath)

	sqlFile, err := os.ReadFile(sqlMigrationFilePath)
	if err != nil {
		log.Fatalf("error migrating db: %v", err)
	}
	sql := string(sqlFile)
	log.Println(string(sqlFile))

	_, err = pgConn.Exec(ctx, sql)
	if err != nil {
		log.Fatalf("error migrating db: %v", err)
	}
}

func DropDb() {
	ctx := context.Background()
	pgConn, err := repository.MakePGConnectionWithUri(PgContainer.URI)
	if err != nil {
		log.Fatalf("error dropping db: %v", err)
	}
	defer pgConn.Close(ctx)
	dropSchemaSQL := "DROP SCHEMA IF EXISTS invoices_service CASCADE"
	_, err = pgConn.Exec(ctx, dropSchemaSQL)
	if err != nil {
		log.Fatalf("error dropping schema: %v", err)
	}
	log.Println("Schema 'invoices_service' dropped successfully")
}
