module store

go 1.17

replace loggers => ./loggers

replace stores => ./stores

replace pkghttp => ./http

require (
	loggers v0.0.0-00010101000000-000000000000
	stores v0.0.0-00010101000000-000000000000
)

require pkghttp v0.0.0-00010101000000-000000000000 // indirect
