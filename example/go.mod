module example

go 1.13

require (
	github.com/GeertJohan/go.rice v1.0.0
	github.com/andrewstucki/web-app-tools/go v0.0.0-20200409181521-2ddaa5e62d69
	github.com/friendsofgo/errors v0.9.2
	github.com/go-chi/chi v4.1.0+incompatible
	github.com/rs/zerolog v1.18.0
	github.com/satori/go.uuid v1.2.0
	github.com/volatiletech/inflect v0.0.0-20170731032912-e7201282ae8d // indirect
	github.com/volatiletech/sqlboiler v3.7.0+incompatible
)

replace github.com/andrewstucki/web-app-tools/go => ../go