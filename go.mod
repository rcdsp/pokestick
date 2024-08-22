module github.com/rcdsp/pokestick

go 1.22.5

require github.com/pelletier/go-toml v1.9.5

require internal/strcase v1.0.0

require (
	github.com/fatih/color v1.17.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/ohler55/ojg v1.24.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
)

replace internal/strcase => ./internal/strcase
