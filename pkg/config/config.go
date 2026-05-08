package config

var Cfg Config

type Config struct {
	System System `mapstructure:"system"`
	MySql  MySql  `mapstructure:"mysql"`
	Redis  Redis  `mapstructure:"redis"`
	JWT    JWT    `mapstructure:"jwt"`
	Email  Email  `mapstructure:"email"`
}

type System struct {
	Env       string `mapstructure:"env"`
	Addr      string `mapstructure:"addr"`
	Location  string `mapstructure:"location"`
	StartCron bool   `mapstructure:"start-cron"`
	ShowSQL   bool   `mapstructure:"show-sql"`
}

type MySql struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Dbname   string `mapstructure:"dbname"`
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
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}
