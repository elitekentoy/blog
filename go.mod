module github.com/elitekentoy/blog

go 1.24.1

replace github.com/elitekentoy/blog/internal/config => ./internal/config
replace "github.com/elitekentoy/blog/internal/database" => ./internal/database
replace "github.com/elitekentoy/blog/internal/api/rss" => ./internal/api/rss

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
)
