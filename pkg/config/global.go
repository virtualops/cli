package config

var GlobalConfig *GlobalConfiguration

type GlobalConfiguration struct {
	AuthToken string `json:"authToken"`
}
