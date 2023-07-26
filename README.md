##

To install sqlc (github.com/kyleconroy/sqlc/cmd/sqlc@latest) you need: `export CGO_ENABLED=1`
After install set it back to 0 `export CGO_ENABLED=0`

To verify: `go env CGO_ENABLED`

## 

for local testing have docker installed

PORT=4000
DB_URL=postgres://iulian:admin@localhost:5432/rssagg?sslmode=disable
