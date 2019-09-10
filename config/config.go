package config

type Config struct {
	Host     string
	Port     string
	DbUser   string
	DbPass   string
	DbName   string
	HttpPort string
	SrvKey   string 
}

func SetConfig() Config {
	var config Config

	//set configuration here
	config.Host = "cockroachdb"
	config.Port = "3306"
	config.DbUser = "root"
	config.DbName = "provisioning_owncloud"
	config.DbPass = ""
	config.HttpPort = "1234"
	config.SrvKey = "Aw4s_g4l4k"
	return config
}
