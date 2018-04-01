package structs

type Rules struct {
	Operations []Operation `json:"enabled"`
}

type Operation struct {
	Description string                 `json:"description"`
	Method      string                 `json:"method"`
	Path        string                 `json:"path"`
	Params      map[string]interface{} `json:"params"`
	Body        map[string]interface{} `json:"body"`
	ParamsCode  string
	BodyCode    string
}

type HostProperties struct {
	Items []HostProperty `json:"properties"`
}

type HostProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Configuration struct {
	Server ServerConfig `json:"server"`
	Host   HostConfig   `json:"host"`
}

type ServerConfig struct {
	Host         string `json:"host" env:"SERVERHOST"`
	Port         int    `json:"port" env:"SERVERPORT"`
	RequestRetry int    `json:"requestRetry" env:"REQUESTRETRY"`
}

type HostConfig struct {
	Destination       string `json:"destination" env:"DESTINATION"`
	UrlPrefix         string `json:"urlPrefix" env:"URLPREFIX"`
	HostRulesURL      string `json:"hostRulesUrl" env:"HOSTRULESURL"`
	HostPropertiesURL string `json:"hostPropertiesUrl" env:"HOSTPROPERTIESURL"`
}
