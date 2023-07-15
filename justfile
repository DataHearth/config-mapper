set dotenv-load
set shell := ["zsh", "-uc"]

default:
  @just --list

publish version:
  sd $(git describe --tags --abbrev=0) {{version}} cmd/cli.go
  git-chglog --next-tag {{version}} --output CHANGELOG.md
  git add CHANGELOG.md cmd/cli.go && git commit -m "chore(changelog): release {{version}}"
  git tag -a {{version}} -m "{{version}}"
  git push --follow-tags