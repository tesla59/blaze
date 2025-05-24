package types

type (
	Config struct {
		Server      *Server
		Db          *Database
		Environment string
	}

	Server struct {
		Host           string
		Port           string
		SSL            *SSL
		Secret         string   `mapstructure:"secret_key"`
		AllowedOrigins []string `mapstructure:"allowed_origins"`
	}

	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Dbname   string
		SSLMode  string
		TimeZone string
	}

	SSL struct {
		Enabled  bool
		CertFile string
		KeyFile  string
	}
)
