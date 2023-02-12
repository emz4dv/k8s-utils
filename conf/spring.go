package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type springCloudConfig struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           string           `json:"label"`
	Version         string           `json:"version"`
	PropertySources []propertySource `json:"propertysources"`
}

type propertySource struct {
	Name   string                 `json:"name"`
	Source map[string]interface{} `json:"source"`
}

func fetchConfiguration(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
			panic("couldn't load configuration, cannot start. terminating. error: " + err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func parseConfiguration(body []byte, config *Config) {

	var cloudconfig springCloudConfig
	err := json.Unmarshal(body, &cloudconfig)

	if err != nil {
		log.Fatal("cannot parse configuration, message: ", err)
	}

	for key, value := range cloudconfig.PropertySources[0].Source {
		if isAuthKey(key) {
			valueString := fmt.Sprintf("%v", value)
			switch key {
				case "iam.baseUri":
					config.Auth.URL = strings.TrimRight(valueString, "auth")
				case "iam-client.username":
					config.Auth.User = valueString	
				case "iam-client.password":
					config.Auth.Password = valueString
				case "iam-client.clientId":
					config.Auth.ClientID = valueString
				case "iam-client.clientSecret":
					config.Auth.ClientSecret = valueString
				}
			//fmt.Printf("loading config property %v => %v\n", key, valueString)
		}
	}
}

func isAuthKey(key string) bool {
    switch key {
    case
        "iam.baseUri",
        "iam-client.username",
        "iam-client.password",
        "iam-client.clientId",
		"iam-client.clientSecret":
        return true
    }
    return false
}

func GetConfigFromConfApp(config *Config) {
	for _, profile :=  range strings.Split(config.Conf.Profiles, ","){
		//fmt.Printf("Parsing %s profile...\n", profile)
		url := fmt.Sprintf("%s/application/%s", config.Conf.URL, profile)
		body, err := fetchConfiguration(url)
		if err!= nil {
			log.Fatal("error:", err)
		}
		parseConfiguration(body, config)
	}
}