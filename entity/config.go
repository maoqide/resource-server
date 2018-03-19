package entity

type MgoConfig struct {
	Addrs    []string `toml:"addrs"` //ip:port
	Database string   `toml:"database"`
	Username string   `toml:"username"`
	Password string   `toml:"password"`
	Timeout  int      `toml:"timeout"`
}
