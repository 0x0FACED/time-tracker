package configs

type DatabaseConfig struct {
	DBUsername string
	DBName     string
	DBHost     string
	DBPort     string
	DBPass     string
}

type ServerConfig struct {
	Host string
	Port string
}

type Config struct {
	DatabaseConfig
	ServerConfig
}
}
