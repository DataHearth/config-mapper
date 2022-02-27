# config-mapper

`config-mapper` is CLI utility tool to help you manage your configuration between systems.  
It provides a set of tools to load your configuration from a system, save it on a git repository and then save it to a new system. This configuration can be a set of files, folders or even dependencies.

## Usage

Before going any further, you need to create a repository to store your configuration. You can choose any supplier as long it's a git repository :).

When copying a file from your configuration repository to your system, it's performing a copy. If the file exists on the system, it's content will be replaced by your configuration's one.

The system is detected automatically. You just need to specify whether the related field in case of `files` or folders `sections` (fields: `darwin` | `linux`).

You can get a configuration template [here](https://raw.githubusercontent.com/DataHearth/config-mapper/main/.config-mapper.yml.template).

### Setup

Create a file called `.config-mapper.yml` in your `home` directory (it is the default search path for config-mapper).  
If you wish to move it to another directory, you can have to choice to inform the tool. By either set an environment like this one: `CONFIG_MAPPER_CFG=/path/to/config/.config-mapper.yml`. Or by providing the `-c /path/to/config/.config-mapper.yml` flag to the tool.

Once the configuratio file created, run this command to initialize the repository localy:

```bash
config-mapper init
```

If the folder is already present and is a git directory, clone instruction will be ignored.

template for storage part:

```yaml
storage:
  # Where will be the repository folder located ? [DEFAULT: MacOS($TMPDIR/config-mapper) | Linux(/tmp/config-mapper)]
  location: /path/to/folder
  git:
    # * by default, if ssh dict is set with its keys filled, I'll try to clone with SSH
    repository: git@github.com:DataHearth/my-config.git
    basic-auth:
      username: USERNAME
      # * NOTE: if you're having trouble with error "authentication required", you should maybe use a token access
      # * In some cases, it's due to 2FA authentication enabled on the git hosting provided
      password: TOKEN
    ssh:
      # path can be relative and can contain environment variables
      private-key: /path/to/private/key
      passphrase: PASSPHRASE
```

### Save your configuration into your repository

Now that your repository is setup localy, you can sync your configuration into it by simply running this command:

```bash
config-mapper save
```

All defined files and folders will be copied inside your repository.

If you want to exclude one part of your configuration file (files, folders, package-managers), you can use these flags to ignore them `--disable-files` `--disable-folders` `--disable-pkgs`

If `homebrew` is provided in the `installation-order` (default: `["apt", "homebrew"]`), it will override the `homebrew` field with all user installed packages (`brew leaves --installed-on-request`). The same principle will be implemented with `aptitude`.

template for your configuration:

```yaml
# NOTE: the $LOCATION if refering to the "storage.location" path. It'll be replaced automatically
# The left part of ":" is your repository location and right part when it should be on your system
files:
  - darwin: "$LOCATION/macos/.zshrc:~/.zshrc"
    linux: "$LOCATION/linux/.zshrc:~/.zshrc"

folders:
  - darwin: "$LOCATION/macos/.config:~/.config"
    linux: "$LOCATION/macos/.config:~/.config"

package-managers:
  installation-order: ["homebrew"]
  homebrew:
    - bat
    - hexyl
    - fd
    - hyperfine
    - diskus
    - jq
    - k9s
    - go
    - starship
    - exa
    - httpie
    - neovim
    - nmap
    - pinentry
    - zsh

  apt-get: []
```

### Load your configuration onto the system

Once your repository is populated with your configurations, you can now load them onto a new system by using:

```bash
config-mapper load
```

The same ignore flags are used in the `save` command.

## TO-DO

- [] load configuration though SSH
- [] save configuration though SSH
- add more storage options
  - [] smb storage
  - [] nfs storage
  - [] zip

## Known issues

- GitHub SSH repository url: `ssh: handshake failed: knownhosts: key mismatch`
  Resolved by create a new primary key based on GitHub new GIT SSH standards ([issue](https://github.com/go-git/go-git/issues/411))
- Cloning from GitHub with `https BasicAuth` and 2FA activated: `authentication required`
  Resolved by creating an access token and set it as password in configuration
