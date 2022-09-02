import enum
import os
import sys
import re
import subprocess
from typing import Any, Dict
import requests

GITEA_API = "https://gitea.antoine-langlois.net/api/v1/repos/DataHearth/config-mapper"
NTR = "\033[0m"  # * Neutral
INF = "\033[0;34m"  # * Blue (info)
WRN = "\033[1;33m"  # * Yellow (warning)
ERR = "\033[1;31m"  # * Red (error)
version_regex = re.compile(r"Version: \"v\d*.\d*.\d*\"")


class LogLevel(enum.Enum):
    INFO = "INFO"
    WARNING = "WARNING"
    ERROR = "ERROR"


def log(msg: str, level: LogLevel = LogLevel.INFO):
    color_lvl = (
        INF
        if level == LogLevel.INFO
        else WRN
        if level == LogLevel.WARNING
        else ERR
        if level == LogLevel.ERROR
        else INF
    )
    print(f"{color_lvl}{level.value}{NTR} {msg}")


if __name__ == "__main__":
    release = input("Enter a release version (vX.Y.Z): ")

    log("updating release version in files")
    with open("cmd/cli.go") as f:
        data = version_regex.sub(f'Version: "{release}"', f.read())

    with open("cmd/cli.go", "w") as f:
        f.write(data)

    res = subprocess.run(
        ["git-chglog", "--next-tag", release, "--output", "CHANGELOG.md"],
        stderr=subprocess.PIPE,
        stdout=subprocess.DEVNULL,
    )
    if res.returncode != 0:
        log(
            f'failed to generate changelog: {res.stderr.decode("UTF-8")}',
            LogLevel.ERROR,
        )
        sys.exit(1)

    log("commit & push changes")
    res = subprocess.run(
        args=f"git add . && git commit -m {release}",
        stderr=subprocess.PIPE,
        stdout=subprocess.DEVNULL,
        shell=True,
    )
    if res.returncode != 0:
        log(
            f'failed to commit changes: {res.stderr.decode("UTF-8")}',
            LogLevel.ERROR,
        )
        sys.exit(1)
    res = subprocess.run(
        args=f"git tag -a {release} -m {release} && git push --follow-tags",
        stderr=subprocess.PIPE,
        stdout=subprocess.DEVNULL,
        shell=True,
    )
    if res.returncode != 0:
        log(
            f'failed to tag and push changes: {res.stderr.decode("UTF-8")}',
            LogLevel.ERROR,
        )
        sys.exit(1)

    log("building Linux binary")
    res = subprocess.run(
        args=["go", "build", "-o", "build/x86-x64_linux_config-mapper"],
        env=os.environ | {"GOOS": "linux"},
        stderr=subprocess.PIPE,
        stdout=subprocess.DEVNULL,
    )
    if res.returncode != 0:
        log(
            f'failed to build linux binary: {res.stderr.decode("UTF-8")}',
            LogLevel.ERROR,
        )
        sys.exit(1)

    log("building Darwin binary")
    res = subprocess.run(
        args=["go", "build", "-o", "build/x86-x64_darwin_config-mapper"],
        env=os.environ | {"GOOS": "darwin"},
        stderr=subprocess.PIPE,
        stdout=subprocess.DEVNULL,
    )
    if res.returncode != 0:
        log(
            f'failed to build darwin binary: {res.stderr.decode("UTF-8")}',
            LogLevel.ERROR,
        )
        sys.exit(1)

    log("creating gitea release")
    api_token: str
    if len(sys.argv) > 1:
        api_token = sys.argv.pop()
    elif "GIT_CFG_MAPPER_TOKEN" in os.environ:
        api_token = os.environ["GIT_CFG_MAPPER_TOKEN"]
    else:
        log("no gitea api token found in CLI params nor in ENV", LogLevel.ERROR)
        sys.exit(1)

    res = subprocess.run(
        args=["git-chglog", "-t", ".chglog/RELEASE_CHANGELOG.tpl.md"],
        stderr=subprocess.PIPE,
        stdout=subprocess.PIPE,
    )
    if res.returncode != 0:
        log(
            f'failed to generate release body: {res.stderr.decode("UTF-8")}',
            LogLevel.ERROR,
        )
        sys.exit(1)

    response = requests.post(
        url=f"{GITEA_API}/releases",
        headers={"Authorization": f"token {api_token}"},
        json={
            "body": res.stdout.decode("UTF-8"),
            "draft": False,
            "prerelease": False,
            "name": release,
            "tag_name": release,
        },
    )
    if not response.ok:
        log(
            f"failed to generate release (status {response.status_code}): {response.json()}",
            LogLevel.ERROR,
        )
        sys.exit(1)

    body: Dict[str, Any] = response.json()
    release_id = body.get("id")
    if not release_id:
        log("no release id found in response body", LogLevel.ERROR)
        sys.exit(1)

    response = requests.post(
        url=f"{GITEA_API}/releases/{release_id}/assets",
        headers={"Authorization": f"token {api_token}"},
        files={
            "attachment": (
                "x86-x64_linux_config-mapper",
                open("build/x86-x64_linux_config-mapper", "rb"),
            )
        },
    )

    if not response.ok:
        log(
            f"failed to upload linux binary (status: {response.status_code}): {response.json()}",
            LogLevel.ERROR,
        )
        sys.exit(1)
    response = requests.post(
        url=f"{GITEA_API}/releases/{release_id}/assets",
        headers={"Authorization": f"token {api_token}"},
        files={
            "attachment": (
                "x86-x64_darwin_config-mapper",
                open("build/x86-x64_darwin_config-mapper", "rb"),
            )
        },
    )
    if not response.ok:
        log(
            f"failed to upload darwin binary (status: {response.status_code}): {response.json()}",
            LogLevel.ERROR,
        )
        sys.exit(1)

    log("Done !")
