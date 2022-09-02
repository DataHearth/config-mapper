package configuration

type Configuration struct {
	Storage         Storage      `mapstructure:"storage" yaml:"storage"`
	Files           []OSLocation `mapstructure:"files" yaml:"files"`
	Folders         []OSLocation `mapstructure:"folders" yaml:"folders"`
	PackageManagers PkgManagers  `mapstructure:"package-managers" yaml:"package-managers"`
}

type OSLocation struct {
	Darwin string `mapstructure:"darwin" yaml:"darwin"`
	Linux  string `mapstructure:"linux" yaml:"linux"`
}

type Storage struct {
	Path string `mapstructure:"location" yaml:"location"`
	Git  Git    `mapstructure:"git" yaml:"git"`
}

type Git struct {
	URL       string    `mapstructure:"repository" yaml:"repository"`
	Name      string    `mapstructure:"name" yaml:"name"`
	Email     string    `mapstructure:"email" yaml:"email"`
	BasicAuth BasicAuth `mapstructure:"basic-auth" yaml:"basic-auth"`
	SSH       interface{}       `mapstructure:"ssh" yaml:"ssh"`
}

type BasicAuth struct {
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
}

type Ssh struct {
	PrivateKey string `mapstructure:"private-key" yaml:"private-key"`
	Passphrase string `mapstructure:"passphrase" yaml:"passphrase"`
}

type PkgManagers struct {
	InstallationOrder []string `mapstructure:"installation-order" yaml:"installation-order"`
	Brew              []string `mapstructure:"brew" yaml:"brew"`
	Apt               []string `mapstructure:"apt" yaml:"apt"`
	Cargo             []string `mapstructure:"cargo" yaml:"cargo"`
	Pip               []string `mapstructure:"pip" yaml:"pip"`
	Npm               []string `mapstructure:"npm" yaml:"npm"`
	Go                []string `mapstructure:"go" yaml:"go"`
	Nala              []string `mapstructure:"nala" yaml:"nala"`
}
