package mapper

type Configuration struct {
	Storage         Storage     `yaml:"storage"`
	Files           []string    `yaml:"files"`
	Folders         []string    `yaml:"folders"`
	PackageManagers PkgManagers `yaml:"package-managers"`
}

type Storage struct {
	Location string `yaml:"location"`
	Git      struct {
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		Repository string `yaml:"repository"`
	} `yaml:"git"`
}

type PkgManagers struct {
	InstallationOrder []string    `yaml:"installation-order"`
	Homebrew          []string    `yaml:"homebrew"`
	AptGet            interface{} `yaml:"apt-get"`
}
