package persistence

var Providers = make(map[string]func(*ConnectionDetails) provider)

type Driver string

type provider interface {
	Driver() Driver
	URL() string
}
