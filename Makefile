.PHONY: migrate-up1 migrate-up2 migrate-up3 migrate-down migrate-create klickhouse-migrate-up db-shell


include .env
export

# Міграції  /schema
migrate-up1:
	docker run --rm \
		-v /home/demian/goshka/first/ticketproject/ticket-service/migrations:/migrations \
		--network ticketproject_ticket-network \
		migrate/migrate \
		-path /migrations \
		-database 'postgres://$(DB_USER):$(DB_PASSWORD)@supportflow-postgres:5432/$(DB_NAME)?sslmode=disable&x-migrations-table=ticket_schema_migrations' \
		up

migrate-up2:
	docker run --rm \
		-v /home/demian/goshka/first/ticketproject/auth-service/migrations:/migrations \
		--network ticketproject_ticket-network \
		migrate/migrate \
		-path /migrations \
		-database 'postgres://$(DB_USER):$(DB_PASSWORD)@supportflow-postgres:5432/$(DB_NAME)?sslmode=disable&x-migrations-table=auth_schema_migrations' \
		up
		
migrate-up3:
	docker exec -i supportflow-clickhouse clickhouse-client \
	--user default \
	--password clickhouse_password \
	--database default < ticket-service/migrations/clickhouse/init.sql

migrate-create:
	docker run --rm \
		-v /home/demian/goshka/first/ticketproject/auth-service/migrations:/migrations \
		migrate/migrate \
		create -ext sql -dir /migrations -seq name=users_table