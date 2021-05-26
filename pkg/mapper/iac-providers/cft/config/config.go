package config

// Config holds the common resource config fields
type Config struct {
	Tags interface{} `json:"tags"`
	Name string      `json:"name"`
}

// AWSResourceConfig helps define type and name for sub-resources if nedded
type AWSResourceConfig struct {
	Resource interface{}
	Name     string
	Type     string
}
