package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// Настройки микросервиса
type Options struct {
	HTTPServer HTTPServer
}

// Инициализация конфигов
func Init(configPath string) (options *Options, err error) {
	if _, err = toml.DecodeFile(configPath, &options); err != nil {
		return nil, errors.Wrap(err, "не удалось загрузить конфиги микросервиса")
	}

	return
}

// Настройки HTTP сервера
type HTTPServer struct {
	Host string
	Port int
}

// Возвращает домен сервера
func (s HTTPServer) GetDomain() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
