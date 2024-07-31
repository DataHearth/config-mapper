# CHANGELOG
<a name="unreleased"></a>
## [Unreleased]


<a name="v0.6.0"></a>
## [v0.6.0] - 2022-08-20
### Bug Fixes
- **configuration:** remove installation-order default value
- **items:** fix stdout when no path is available
- **pkgs:** update command building
- **pkgs:** add pkg manager validation check and parsing cli arguments

### Features
- **packages:** add nala package manager


<a name="v0.5.0"></a>
## [v0.5.0] - 2022-08-01
### Bug Fixes
- **git:** use go-git for adding removed file (workaround)

### Features
- **pkgs:** add more package manager


<a name="v0.4.0"></a>
## [v0.4.0] - 2022-08-01
### Bug Fixes
- **config:** don't throw error when file not available on OS
- **save:** remove folder before copy (avoid unwanted files)

### Features
- **cli:** add verbose flag and a spinner for pkgs
- **config:** add SSH capability with user/pass or key/pass


<a name="v0.3.0"></a>
## [v0.3.0] - 2022-08-01
### Features
- **cli:** packages are disabled by default
- **sync:** add .ignore file to filter folder's content


<a name="v0.2.0"></a>
## [v0.2.0] - 2022-08-01
### Bug Fixes
- **git:** use git binary for "git add"
- **git:** deleted files are not pushed
- **git:** add error handling and repo URL from config
- **rm:** remove dir even if not empty
- **stderr:** add stderr to brew command output

### Code Refactoring
- **archi:** reduce base code to one struct
- **cli:** separate functions from CLI for lisibility
- **logging:** drop pterm

### Features
- **cli:** add git push option with message
- **index:** add indexing system


<a name="v0.1.0"></a>
## v0.1.0 - 2022-07-31
### Bug Fixes
- **config:** fix config path check
- **copy:** use io.Copy instead of custom copy

### Code Refactoring
- **config:** unmarshal configuration instead of raw read

### Features
- **cli:** add save and load features
- **cli:** add init sub-command
- **cli:** add copy folder
- **cli:** add save command
- **cli:** add configuration-file persistant flag
- **cli:** implement pkgs installation
- **config:** update git configuration
- **config:** add yaml tags for yaml.v3


[Unreleased]: https://gitea.antoine-langlois.net/DataHearth/config-mapper/compare/v0.6.0...HEAD
[v0.6.0]: https://gitea.antoine-langlois.net/DataHearth/config-mapper/compare/v0.5.0...v0.6.0
[v0.5.0]: https://gitea.antoine-langlois.net/DataHearth/config-mapper/compare/v0.4.0...v0.5.0
[v0.4.0]: https://gitea.antoine-langlois.net/DataHearth/config-mapper/compare/v0.3.0...v0.4.0
[v0.3.0]: https://gitea.antoine-langlois.net/DataHearth/config-mapper/compare/v0.2.0...v0.3.0
[v0.2.0]: https://gitea.antoine-langlois.net/DataHearth/config-mapper/compare/v0.1.0...v0.2.0
