module main

go 1.21.5

require (
	api v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.7.1
	github.com/gorilla/mux v1.8.1
	webInteract v0.0.0-00010101000000-000000000000
)

require (
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/sessions v1.2.2 // indirect
	golang.org/x/crypto v0.17.0 // indirect
)

replace api => ./api

replace webInteract => ./webInteract
