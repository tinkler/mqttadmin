package conf

type Conf struct {
	Db     *DbConfig
	Server *ServerConfig
}

func NewConf() *Conf {
	return &Conf{Server: newServerConfig(), Db: &DbConfig{Dsn: "host=localhost user=mcrz password=mcrz4105 dbname=mcrz port=5432 sslmode=disable TimeZone=Asia/Shanghai"}}
}
