package types

type (
	Config struct {
		Server *Server
		Db     *Database
	}

	Server struct {
		Host string
		Port string
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
)
