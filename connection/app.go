package connection

import (
	//"encoding/json"
	//"fmt"
	"encoding/json"
	"log"

	"github.com/go-resty/resty/v2"
)
type ActuatorHealth struct {
	Status     string `json:"status"`
	Components struct {
		ClientConfigServer struct {
			Status  string `json:"status"`
			Details struct {
				PropertySources []string `json:"propertySources"`
			} `json:"details"`
		} `json:"clientConfigServer"`
		DiscoveryComposite struct {
			Description string `json:"description"`
			Status      string `json:"status"`
			Components  struct {
				DiscoveryClient struct {
					Description string `json:"description"`
					Status      string `json:"status"`
				} `json:"discoveryClient"`
			} `json:"components"`
		} `json:"discoveryComposite"`
		DiskSpace struct {
			Status  string `json:"status"`
			Details struct {
				Total     int64 `json:"total"`
				Free      int64 `json:"free"`
				Threshold int   `json:"threshold"`
				Exists    bool  `json:"exists"`
			} `json:"details"`
		} `json:"diskSpace"`
		LivenessState struct {
			Status string `json:"status"`
		} `json:"livenessState"`
		Ping struct {
			Status string `json:"status"`
		} `json:"ping"`
		ReactiveDiscoveryClients struct {
			Description string `json:"description"`
			Status      string `json:"status"`
			Components  struct {
				SimpleReactiveDiscoveryClient struct {
					Description string `json:"description"`
					Status      string `json:"status"`
				} `json:"Simple Reactive Discovery Client"`
			} `json:"components"`
		} `json:"reactiveDiscoveryClients"`
		ReadinessState struct {
			Status string `json:"status"`
		} `json:"readinessState"`
		RefreshScope struct {
			Status string `json:"status"`
		} `json:"refreshScope"`
	} `json:"components"`
	Groups []string `json:"groups"`
}

type Swagger struct {
	Swagger string `yaml:"swagger"`
	Info    struct {
	} `yaml:"info"`
	Host     string        `yaml:"host"`
	BasePath string        `yaml:"basePath"`
	Tags     []interface{} `yaml:"tags"`
	Paths    map[string]interface{} `yaml:"paths"`
	Definitions struct {
	} `yaml:"definitions"`
}

func getAppJson(url string, token string) []byte{
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		SetHeader("Accept", "application/json").
		SetAuthToken(token).
		Get(url)
	if err != nil {
		log.Fatal("error:", err)
	} 

	return resp.Body()
}


func GetActuatorHealth(url string, token string) ActuatorHealth{
	var result ActuatorHealth

	err := json.Unmarshal(getAppJson(url, token), &result)

	if err != nil {
		log.Fatal("cannot parse configuration, message: ", err)
	}

	return result
}

//Временно, в будущем необходимо использовать generics
func GetSwagger(url string, token string) Swagger{
	var result Swagger

	err := json.Unmarshal(getAppJson(url, token), &result)

	if err != nil {
		log.Fatal("cannot parse configuration, message: ", err)
	}

	return result
}