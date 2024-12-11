package config

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Host string
	Port string
}

func (c *Config) New() {
	config := ParseConfig()
	c.Host = config.Host
	c.Port = config.Port
}
func ParseConfig() *Config {
	yamlFile, err := ioutil.ReadFile("/home/midas/GolandProjects/PaymentSystem/config/config.yaml")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to read config file")
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to parse config file")
	}
	zap.S().Info("Loaded config.yaml")
	return &config
}
