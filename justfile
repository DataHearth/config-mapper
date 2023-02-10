set shell := ["zsh", "-uc"]

default:
  @just --list

publish version:
  git-chglog --next-tag {{version}} --output CHANGELOG.md
  git add CHANGELOG.md && git commit -m "chore(changelog): release {{version}}"
  git tag -a {{version}} -m "{{version}}"
  git push --follow-tags