latest-tag := `git describe --tags --abbrev=0`

publish version: (bump-files-version version)
  git-chglog --next-tag {{version}} --output CHANGELOG.md
  git add CHANGELOG.md cmd/cli.go && git commit -m "chore(changelog): release {{version}}"
  git tag -a {{version}} -m "{{version}}"
  git push --follow-tags

bump-files-version version:
  sd {{latest-tag}} {{version}} cmd/cli.go
  sd {{latest-tag}} "version-{{version}}-blue" CHANGELOG.md