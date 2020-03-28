package configure

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

type Config struct {
	Address string `json:"localaddress"`
	Dns struct{
		Bind string `json:"bind"`
		Port int `json:"port"`
	} `json:"dns"`
	GRpc struct{
		Address string `json:"address"`	
		Port int `json:"port"`
	} `json:"grpc"`
}

var (
	configInstance *Config
	once sync.Once
)

func init() {
	viper.AddConfigPath("./")
	viper.AddConfigPath("configure/")
}



func NewConfig()*Config{
	once.Do(func() {
		configInstance = &Config{}
		configInstance.readFromFile()
	})
	return configInstance
}

func(c *Config)readFromFile(){
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


