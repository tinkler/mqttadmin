package conf

type Conf struct {
	Db     *DbConfig
	Server *ServerConfig
}

func NewConf() *Conf {
	return &Conf{Server: newServerConfig(), Db: newDbConfig()}
}
