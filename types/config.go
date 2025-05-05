package types

type (
	Config struct {
		Server *Server
		Db     *Database
	}

	Server struct {
		Host string
		Port string
		SSL  *SSL
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
