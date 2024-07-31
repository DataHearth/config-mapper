#!/bin/bash

set -e

VERSION=v0.5.0

log() {
    NTR=$'\033[0m'    # * Neutral
    INF=$'\033[0;34m' # * Blue (info)
    WRN=$'\033[1;33m' # * Yellow (warning)
    ERR=$'\033[1;31m' # * Red (error)

    log_lvl=""
    case $1 in
        INFO)
            log_lvl="${INF}$1"
            ;;
        WARNING)
            log_lvl="${WRN}$1"
            ;;
        ERROR)
            log_lvl="${ERR}$1"
            ;;
    esac

    log_lvl="${log_lvl}${NTR}"
    msg="${log_lvl}\t$2"

    echo -e "${msg}"
}

log "INFO" "checking required dependencies to create release"
if ! type git 1> /dev/null; then
  log "ERROR" "\"git\" binary not available"
  exit 1
fi
if ! type sd 1> /dev/null; then
  log "ERROR" "\"sd\" binary not available"
  exit 1
fi
if ! type gh 1> /dev/null; then
  log "ERROR" "\"gh\" binary not available"
  exit 1
fi
if ! type go 1> /dev/null; then
  log "ERROR" "\"go\" binary not available"
  exit 1
fi
if ! type git-chglog 1> /dev/null; then
  log "ERROR" "\"git-chglog\" binary not available"
  exit 1
fi
if ! type jq 1> /dev/null; then
  log "ERROR" "\"jq\" binary not available"
  exit 1
fi
if ! type xhs 1> /dev/null; then
  log "ERROR" "\"xhs\" binary not available"
  exit 1
fi

read -p "Enter a release version (vX.Y.Z): " release

log "INFO" "updating release version in files"
sd "Version: \"$VERSION\"" "Version: \"$release\"" cmd/cli.go
sd "VERSION=$VERSION" "VERSION=$release" release.sh

log "INFO" "updating changelog"
git-chglog --next-tag $release --output CHANGELOG.md

log "INFO" "commit & push changes"
git add .
git commit -m "$release"
git push
git tag -a $release -m $release
git push --tags

log "INFO" "building Linux binary"
GOOS=linux go build -o build/x86-x64_linux_config-mapper

log "INFO" "building Darwin binary"
GOOS=darwin go build -o build/x86-x64_darwin_config-mapper

log "INFO" "creating release"
local response=$(xhs POST https://gitea.antoine-langlois.net/api/v1/repos/DataHearth/config-mapper/releases Authorization:"token $GIT_CFG_MAPPER_TOKEN"  body=$(git-chglog -t .chglog/RELEASE_CHANGELOG.tpl.md) draft:=false name=$release prerelease:=false tag_name=$release)
local release_id=$(echo $reponse | jq .id)

xhs POST https://gitea.antoine-langlois.net/api/v1/repos/DataHearth/config-mapper/releases/$release_id/assets name=="x86-x64_linux_config-mapper" Authorization:"token $GIT_CFG_MAPPER_TOKEN" attachement@build/x86-x64_linux_config-mapper
xhs POST https://gitea.antoine-langlois.net/api/v1/repos/DataHearth/config-mapper/releases/$release_id/assets name=="x86-x64_darwin_config-mapper" Authorization:"token $GIT_CFG_MAPPER_TOKEN" attachement@build/x86-x64_darwin_config-mapper
