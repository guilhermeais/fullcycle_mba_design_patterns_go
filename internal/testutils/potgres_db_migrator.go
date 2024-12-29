package testutils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type PostgresDbMigrator struct {
	Conn   pgx.Conn
	Schema string
}

func (p *PostgresDbMigrator) MigrateDb() {
	ctx := context.Background()
	sqlMigrationFilePath := os.Getenv("TEST_SQL_MIGRATION_PATH")

	log.Printf("migrating with file '%s'", sqlMigrationFilePath)

	sqlFile, err := os.ReadFile(sqlMigrationFilePath)
	if err != nil {
		log.Fatalf("error migrating db: %v", err)
	}
	sql := string(sqlFile)
	log.Println(string(sqlFile))

	_, err = p.Conn.Exec(ctx, sql)
	if err != nil {
		log.Fatalf("error migrating db: %v", err)
	}
}

func (p *PostgresDbMigrator) DropDb() {
	ctx := context.Background()
	dropSchemaSQL := fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", p.Schema)
	_, err := p.Conn.Exec(ctx, dropSchemaSQL)
	if err != nil {
		log.Fatalf("error dropping schema '%s': %v", p.Schema, err)
	}
	log.Printf("Schema '%s' dropped successfully\n", p.Schema)
}
