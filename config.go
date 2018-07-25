package gostore

import (
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"gopkg.in/go-playground/validator.v9"
)

type GSConfig struct {
	SavePath  string `validate: "required"`
	CrawlSize int    `validate: "required"`
}

func GetConfig(path string) (GSConfig, error) {
	var conf GSConfig

	config.Load(file.NewSource(file.WithPath(path)))
	config.Scan(&conf)

	vd := validator.New()
	err := vd.Struct(conf)
	if err != nil {
		return GSConfig{}, err
	}
	return conf, err
}
