package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

//FileName где лежит.
const FileName = "configs/config.toml"

var sc *ServerConfig

//ServerConfig - all cfg data
type ServerConfig struct {
	DB  Database `toml:"database"`
	Srv Server   `toml:"server"`
}

//Database config
type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	NameDB   string
}

//Server config
type Server struct {
	Host string
	Port string
}

//Load from file
func Load(filename string) (err error) {
	sc = new(ServerConfig)
	_, err = toml.DecodeFile(filename, sc)

	return err
}

//GetServerAddress возвращает адрес, на котором должен запуститься сервер.
func GetServerAddress() (result string) {
	result = ":3333" //по-умолчанию (для локального запуска)
	if nil == sc {
		return
	}
	cfg := sc.Srv
	if "" == cfg.Host && "" == cfg.Port {
		return
	}
	result = fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	return
}

//GetInfoDB возвращает инфо о БД.
func GetInfoDB() (result string) {
	result = "(нет)"
	if nil == sc {
		return
	}
	cfg := sc.DB
	result = fmt.Sprintf("%s:%s (%s)", cfg.Host, cfg.Port, cfg.NameDB)
	return
}

//GetConnectionStringDB возвращает строку для соединения с БД.
func GetConnectionStringDB() (result string) {
	if nil == sc {
		return
	}
	cfg := sc.DB

	result = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.NameDB)

	return
}
