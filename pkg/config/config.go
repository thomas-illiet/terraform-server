package config

// Server defines the server configuration.
type Server struct {
	Addr          string `json:"addr"`
	Cert          string `json:"cert"`
	Key           string `json:"key"`
	StrictCurves  bool   `json:"strict_curves"`
	StrictCiphers bool   `json:"strict_ciphers"`
}

// Metrics defines the metrics server configuration.
type Metrics struct {
	Addr  string `json:"addr"`
	Token string `json:"token"`
}

// Logs defines the level and color for log configuration.
type Logs struct {
	Level  string `json:"level"`
	Pretty bool   `json:"pretty"`
	Color  bool   `json:"color"`
}

// Encryption defines the encryption configuration.
type Encryption struct {
	Secret string `json:"secret"`
}

// Config defines the general configuration.
type Config struct {
	Server     Server     `json:"server"`
	Metrics    Metrics    `json:"metrics"`
	Logs       Logs       `json:"log"`
	Encryption Encryption `json:"encryption"`
}

// Load initializes a default configuration struct.
func Load() *Config {
	return &Config{}
}
