DB_URL=postgres://docker:docker@localhost:5432/encurtador_de_urls?sslmode=disable

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" down

migrate-drop:
	migrate -path db/migrations -database "$(DB_URL)" drop

migrate-force:
	migrate -path db/migrations -database "$(DB_URL)" force 0

migrate-version:
	migrate -path db/migrations -database "$(DB_URL)" version
	