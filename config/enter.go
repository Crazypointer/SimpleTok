package config

import "strconv"

type Mysql struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DB       string `yaml:"db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
type Server struct {
	BaseUrl string `yaml:"base_url"`
}
type Cos struct {
	BucketName string `yaml:"bucket_name"`
	CosRegion  string `yaml:"cos_region"`
	SecretID   string `yaml:"secret_id"`
	SecretKey  string `yaml:"secret_key"`
	CosUrl     string `yaml:"cos_url"`
}
type Local struct {
	Enable bool `yaml:"enable"`
}
type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool-size"`
}
type Jwt struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"`
	Issuer string `yaml:"issuer"`
}

type Config struct {
	Mysql  Mysql  `yaml:"mysql"`
	Server Server `yaml:"server"`
	Cos    Cos    `yaml:"cos"`
	Local  Local  `yaml:"local"`
	Redis  Redis  `yaml:"redis"`
	Jwt    Jwt    `yaml:"jwt"`
}

func (m Mysql) Dsn() string {
	return m.User + ":" + m.Password + "@tcp(" + m.Host + ":" + strconv.Itoa(m.Port) + ")/" + m.DB + "?charset=utf8mb4&parseTime=True&loc=Local"
}
