module github.com/KelvinMcClean/palimpsest/hardcover

go 1.26.5

require github.com/pelletier/go-toml v1.9.5

require (
	github.com/KelvinMcClean/palimpsest v0.0.0-20260819173055-9b1ff8ce5eb1 // indirect
	github.com/KelvinMcClean/palimpsest/database v0.0.0-20260814105229-be9ca5289805 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.10.0 // indirect
	golang.org/x/text v0.29.0 // indirect
)

replace github.com/KelvinMcClean/palimpsest/database => ../database
replace github.com/KelvinMcClean/palimpsest => ..
