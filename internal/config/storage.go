package config

type Storage struct {
	Git interface{} `yaml:"git"`
	AWS interface{}
	GCP interface{}
}

type Git struct {
	Repository     string `yaml:"repository"`
	Authentication struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"authentication"`
}
