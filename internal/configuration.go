package mapper

type Configuration struct {
	Storage         Storage     `yaml:"storage"`
	Files           []string    `yaml:"files"`
	Folders         []string    `yaml:"folders"`
	PackageManagers PkgManagers `yaml:"package-managers"`
}

type Storage struct {
	Location string `yaml:"location"`
	Git      Git    `yaml:"git"`
}

type Git struct {
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Repository string `yaml:"repository"`
}

type PkgManagers struct {
	InstallationOrder []string    `yaml:"installation-order"`
	Homebrew          []string    `yaml:"homebrew"`
	AptGet            interface{} `yaml:"apt-get"`
}
