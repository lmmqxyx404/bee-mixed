package config

type BeeShop struct {
	Disable bool   `mapstructure:"disable" json:"disable" yaml:"disable"`
	Listen  string `mapstructure:"listen" json:"listen" yaml:"listen"`
	Host    string `mapstructure:"host" json:"host" yaml:"host"`
}
