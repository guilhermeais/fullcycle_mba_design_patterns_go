package testutils

import (
	"context"
	"fmt"
	"invoices/internal/app/infrastructure/repository"
	"log"
	"net"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	Container testcontainers.Container
	URI       string
}

var PgContainer *PostgresContainer

func init() {
	log.Println("Starting postgres container...")
	port, err := findAvailablePort()
	if err != nil {
		log.Fatalf("error on initiating postgres container: %v", err)
	}
	log.Printf("Starting postgres on port '%d'...\n", port)
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{fmt.Sprintf("%d:5432/tcp", port)},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.AutoRemove = true
		},
	}
	ctx := context.Background()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("error on initiating postgres container: %v", err)
	}
	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("error on initiating postgres container: %v", err)
	}
	uri := fmt.Sprintf("postgres://test:test@%s:%d/testdb?sslmode=disable", host, port)
	pgConn, err := repository.MakePGConnectionWithUri(uri)
	if err != nil {
		log.Fatalf("error migrating db: %v", err)
	}
	createExtensionSQL := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	_, err = pgConn.Exec(ctx, createExtensionSQL)
	if err != nil {
		log.Fatalf("error creating extension: %v", err)
	}
	pgConn.Close(ctx)

	PgContainer = &PostgresContainer{
		Container: container,
		URI:       uri,
	}
}

type newSchemaResult struct {
	URI, Schema string
}

func CreateNewSchema() (newSchemaResult, error) {
	ctx := context.Background()
	conn, err := repository.MakePGConnectionWithUri(PgContainer.URI)
	schema := newSchemaResult{}
	if err != nil {
		return schema, fmt.Errorf("error connecting to database: %v", err)
	}
	defer conn.Close(ctx)

	newSchemaName := "schema_" + strings.Join(strings.Split(uuid.New().String(), "-"), "_")
	log.Printf("creating new schema: '%s'", newSchemaName)
	createSchemaSQL := fmt.Sprintf("CREATE SCHEMA %s", newSchemaName)
	_, err = conn.Exec(ctx, createSchemaSQL)
	if err != nil {
		return schema, fmt.Errorf("error creating new schema: %v", err)
	}

	newSchemaURI := fmt.Sprintf("%s&search_path=%s,public", PgContainer.URI, newSchemaName)
	schema.Schema = newSchemaName
	schema.URI = newSchemaURI
	return schema, nil
}

func findAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, fmt.Errorf("failed to find an available port: %v", err)
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}
