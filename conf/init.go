package conf

import (
	"log"
	"github.com/spf13/viper"
)
	
var configViper = viper.New()
type Config struct {
	Tests struct {
		Springboot struct {
			Since string `yaml:"since"`
			Timeouts struct {
				Counts  int `yaml:"counts"`
				Seconds int `yaml:"seconds"`
			} `yaml:"timeouts"`	
		} `yaml:"springboot"`
		PodStatus struct {
			Timeouts struct {
				Counts  int `yaml:"counts"`
				Seconds int `yaml:"seconds"`
			} `yaml:"timeouts"`	
		} `yaml:"podStatus"`
	} `yaml:"tests"`
	App struct {
		Name string `yaml:"name"`
		Artifact string `yaml:"artifact"`
	} `yaml:"app"`
	Conf struct {
		URL string `yaml:"url"`
		Profiles string `yaml:"profiles"`
	} `yaml:"conf"`
	Auth struct {
		URL          string `yaml:"url"`
		Realm        string `yaml:"realm"`
		ClientID     string `yaml:"clientId"`
		ClientSecret string `yaml:"clientSecret"`
		User         string `yaml:"username"`
		Password     string `yaml:"password"`
	} `yaml:"auth"`
	K8S struct {
		InCluster bool `yaml:"inCluster"`
		Cluster string `yaml:"cluster"`
		Namespace string `yaml:"namespace"`
	} `yaml:"k8s"`
}

func LoadConfig() *Config {

	var config Config

	configViper.AddConfigPath("../")
	//configViper.AddConfigPath(".")
	configViper.SetConfigName("config")
	configViper.SetConfigType("yaml") 	

	if err := configViper.ReadInConfig(); err != nil {
		log.Fatal("error:", err)
	}

	if err := configViper.Unmarshal(&config); err != nil {
		log.Fatal("error:", err)
	}
	GetConfigFromConfApp(&config)

	return &config
}




