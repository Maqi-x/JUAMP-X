module JUAMP-X

go 1.23.3

toolchain go1.23.4

require (
	BetterString v0.0.0-00010101000000-000000000000
	Mqio v0.0.0-00010101000000-000000000000
	github.com/BurntSushi/toml v1.4.0
	golang.org/x/term v0.27.0
)

require golang.org/x/sys v0.28.0 // indirect

replace (
	BetterString => ./Include/BetterString
	Mqio => ./Mqio
)
