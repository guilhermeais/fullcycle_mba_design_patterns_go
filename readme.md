## Estrutura do Projeto
├── cmd/
│   └── myproject/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── application/
│   │   │   ├── usecase/
│   │   │   │   ├── generate_invoices.go
│   │   │   │   ├── generate_invoices_test.go
│   │   │   │   └── ...
│   │   ├── domain/
│   │   │   ├── entities/
│   │   │   │   ├── entity.go
│   │   │   │   ├── entity_test.go
│   │   ├── infrastructure/
│   │   │   ├── repository/
│   │   │   │   ├── repository.go
│   │   │   │   ├── repository_test.go
│   │   │   │   └── ...
│   │   │   ├── http/
│   │   │   │   ├── handler.go
│   │   │   │   ├── handler_test.go
│   │   │   │   └── ...
│   │   │   └── ...
│   │   └── ...
│   └── ...
├── pkg/
│   └── ...
├── sql/
│   ├── 001_create_tables.sql
│   ├── 002_insert_data.sql
│   └── ...
├── go.mod
└── go.sum