package config

type IViperConfig interface {
	GetInt(key string, defaultValue int) int
	GetString(key string, defaultValue string) string
	GetBool(key string, defaultValue bool) bool
}
