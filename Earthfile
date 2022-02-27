VERSION 0.6

FROM golang:1.17-alpine3.15
WORKDIR /config-mapper

build-macos:
  COPY . .
  RUN GOOS=darwin go build -o build/config-mapper main.go
  SAVE ARTIFACT build/config-mapper /config-mapper AS LOCAL build/x86-x64_darwin_config-mapper

build-linux:
  COPY . .
  RUN GOOS=linux go build -o build/config-mapper main.go
  SAVE ARTIFACT build/config-mapper /config-mapper AS LOCAL build/x86-x64_linux_config-mapper