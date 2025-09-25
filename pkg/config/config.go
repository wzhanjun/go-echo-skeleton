package config

var Cfg Config

type Config struct {
	System System `yaml:"system"`
	MySql  MySql  `yaml:"mysql"`
	Redis  Redis  `yaml:"redis"`
	// jwt
	JWT JWT `json:"jwt" yaml:"jwt"`
	// email
	Email Email `json:"email" yaml:"email"`
}

type System struct {
	ENV       string `mapstructure:"env"`
	Addr      string `mapstructure:"addr"`
	Location  string `mapstructure:"location"`
	StartCron bool   `mapstructure:"start-cron"`
}

type MySql struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Dbname   string `json:"dbname" yaml:"dbname"`
}

type Redis struct {
	DB       int         `mapstructure:"db"`
	Addr     string      `mapstructure:"addr"`
	Password string      `mapstructure:"password"`
	Sentinel SentinelCfg `mapstructure:"sentinel"`
}

type SentinelCfg struct {
	MasterName string   `mapstructure:"master-name"`
	Nodes      []string `mapstructure:"nodes"`
	Password   string   `mapstructure:"password"`
}

type JWT struct {
	// issuer
	Issuer string `mapstructure:"issuer"`
	// secret
	Secret string `mapstructure:"secret"`
	// expires
	Expires int `mapstructure:"expires"`
}

type Email struct {
	User string `mapstructure:"user" yaml:"user"`
	Pass string `mapstructure:"pass" yaml:"pass"`
	Host string `mapstructure:"host" yaml:"host"`
	Port string `mapstructure:"port" yaml:"port"`
}
