module example

go 1.13

require (
	github.com/GeertJohan/go.rice v1.0.0
	github.com/andrewstucki/web-app-tools/go v0.0.0-20200409181521-2ddaa5e62d69
	github.com/go-chi/chi v4.1.0+incompatible
	github.com/golang-migrate/migrate/v4 v4.10.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/lib/pq v1.3.0
	github.com/rs/zerolog v1.18.0
)

replace github.com/andrewstucki/web-app-tools/go => ../go
