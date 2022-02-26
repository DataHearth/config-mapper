# config-mapper

## Known issues

- GitHub SSH repository url: `ssh: handshake failed: knownhosts: key mismatch`
  Resolved by create a new primary key based on GitHub new GIT SSH standards ([issue](https://github.com/go-git/go-git/issues/411))
- Cloning from GitHub with `https BasicAuth` and 2FA activated: `authentication required`
  Resolved by creating an access token and set it as password in configuration
