# CHANGELOG
<a name="unreleased"></a>
## [Unreleased]


<a name="v0.3.0"></a>
## [v0.3.0] - 2022-06-01
### Features
- **cli:** packages are disabled by default
- **sync:** add .ignore file to filter folder's content


<a name="v0.2.0"></a>
## [v0.2.0] - 2022-05-23
### Bug Fixes
- **config:** fix config path check
- **copy:** use io.Copy instead of custom copy
- **git:** use git binary for "git add"
- **git:** deleted files are not pushed
- **git:** add error handling and repo URL from config
- **rm:** remove dir even if not empty
- **stderr:** add stderr to brew command output

### Code Refactoring
- **archi:** reduce base code to one struct
- **cli:** separate functions from CLI for lisibility
- **config:** unmarshal configuration instead of raw read
- **logging:** drop pterm

### Features
- **cli:** add configuration-file persistant flag
- **cli:** add git push option with message
- **cli:** add save and load features
- **cli:** add init sub-command
- **cli:** add copy folder
- **cli:** add save command
- **cli:** implement pkgs installation
- **config:** update git configuration
- **config:** add yaml tags for yaml.v3
- **index:** add indexing system


<a name="v0.1.0"></a>
## v0.1.0 - 2022-02-27
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


[Unreleased]: https://github.com/DataHearth/config-mapper/compare/v0.3.0...HEAD
[v0.3.0]: https://github.com/DataHearth/config-mapper/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/DataHearth/config-mapper/compare/v0.1.0...v0.2.0
