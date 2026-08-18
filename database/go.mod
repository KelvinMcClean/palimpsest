module github.com/KelvinMcClean/palimpsest/database

go 1.26.5

require (
	github.com/KelvinMcClean/palimpsest/calibre v0.0.0-20260817210309-d8db81c14ac6
	github.com/KelvinMcClean/palimpsest/hardcover v0.0.0-20260814105229-be9ca5289805
	github.com/jackc/pgx/v5 v5.10.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	golang.org/x/text v0.29.0 // indirect
)

replace github.com/KelvinMcClean/palimpsest/hardcover => ../hardcover

replace github.com/KelvinMcClean/palimpsest/calibre => ../calibre
