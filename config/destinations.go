package config

type Destination interface {
	Name() string
	Folders() []string
	Type() string
	Credentials() map[string]string
}
