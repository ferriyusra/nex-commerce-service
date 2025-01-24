create migration
- migrate create -ext sql -dir database/migrations -seq create_table_account

run migration

migrate -database "postgres://postgres.euuopddxoippbxdovtxo:nex12345@aws-0-ap-southeast-1.pooler.supabase.com:5432/postgres" -path database/migrations up