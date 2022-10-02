set shell := ["zsh", "-uc"]
set dotenv-load

default:
  @just --list

publish version:
  git-chglog --next-tag {{version}} --output CHANGELOG.md
  git add CHANGELOG.md && git commit -m "chore: update CHANGELOG {{version}}"
  git tag -a {{version}} -m "{{version}}"
  git push --follow-tags
  goreleaser release --rm-dist --release-notes <(git-chglog -t .chglog/RELEASE_CHANGELOG.tpl.md)