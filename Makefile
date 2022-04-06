dev:
	go run cmd/apiserver/main.go

migrate_create:
	
setup:
	go mod tidy
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest