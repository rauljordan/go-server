server:
	go build .
mocks:
	mockgen -package=mocks -destination ./internal/mocks/mock_db.go github.com/rauljordan/go-server/internal/db Database
	mockgen -package=mocks -destination ./internal/mocks/mock_server.go github.com/rauljordan/go-server/server Server
migrate:
	migrate -database ${POSTGRESQL_URL} -path ./sql up