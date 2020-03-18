package configure

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Address string `json:"address"`
}

func init() {
	viper.AddConfigPath("./")
	viper.AddConfigPath("configure/")
}



func NewConfig()*Config{
	return &Config{}
}

func(c *Config)ReadFromFile(){
	var filename string = "config"
	viper.SetConfigName(filename)
	if err := viper.ReadInConfig();err != nil {
		logrus.WithError(err).Fatalf("Read Config file %s\n", filename)
	}
	if err := viper.Unmarshal(c);err != nil {
		logrus.WithError(err).Fatalf("Unmarshal data to configure failed\n")
	}
	logrus.Infof(c.Address)
}