create migration
- migrate create -ext sql -dir database/migrations -seq create_table_order

run migration

migrate -database "postgres://postgres.euuopddxoippbxdovtxo:nex12345@aws-0-ap-southeast-1.pooler.supabase.com:5432/postgres" -path database/migrations up

migrate -database "postgres://ferri@127.0.0.1:5432/nex-commerce-service?sslmode-allow" -path database/migrations up