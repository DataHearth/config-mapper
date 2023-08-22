package internal

type Configuration struct {
	Path            string      `yaml:"path"`
	Storage         Storage     `yaml:"storage"`
	Items           []Item      `yaml:"items"`
	PackageManagers PkgManagers `yaml:"package-managers"`
	Logging         Logging     `yaml:"logging"`
}

type Logging struct {
	Level          string `yaml:"log-level,omitempty"`
	EnableTime     bool   `yaml:"time,omitempty"`
	TimeFormat     string `yaml:"time-format,omitempty"`
	File           string `yaml:"file,omitempty"`
	DisableConsole bool   `yaml:"disable-console"`
}

type Item struct {
	Name      string       `yaml:"name"`
	Universal ItemLocation `yaml:"universal"`
	Darwin    ItemLocation `yaml:"darwin"`
	Linux     ItemLocation `yaml:"linux"`
	Windows   ItemLocation `yaml:"windows"`
}

type ItemLocation struct {
	Local  string `yaml:"local"`
	Remote string `yaml:"remote"`
}

type Storage struct {
	Git Git `yaml:"git"`
}

type Git struct {
	Repository string    `yaml:"repository"`
	Name       string    `yaml:"name"`
	Email      string    `yaml:"email"`
	BasicAuth  BasicAuth `yaml:"basic-auth"`
	SSH        []SshAuth `yaml:"ssh-auth"`
}

type BasicAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type SshAuth struct {
	PrivateKey string `yaml:"private-key"`
	Passphrase string `yaml:"passphrase"`
}

type PkgManagers struct {
	InstallationOrder []string `yaml:"installation-order"`
	Brew              []string `yaml:"brew"`
	Apt               []string `yaml:"apt"`
	Cargo             []string `yaml:"cargo"`
	Pip               []string `yaml:"pip"`
	Npm               []string `yaml:"npm"`
	Go                []string `yaml:"go"`
	Nala              []string `yaml:"nala"`
}
