# 🔋 AiRISTA Flow RTLS - BatteryHistory

[![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)
[![Release](https://img.shields.io/badge/Calver-YY.WW.REVISION-22bfda.svg)](https://calver.org/)
[![Linters](https://github.com/Plaenkler/BatteryHistory/actions/workflows/linters.yml/badge.svg)](https://github.com/Plaenkler/BatteryHistory/actions/workflows/linters.yml)
[![CodeQL](https://github.com/Plaenkler/BatteryHistory/actions/workflows/codeql.yml/badge.svg)](https://github.com/Plaenkler/BatteryHistory/actions/workflows/codeql.yml)
[![Goreport](https://goreportcard.com/badge/github.com/Plaenkler/BatteryHistory)](https://goreportcard.com/report/github.com/Plaenkler/BatteryHistory)

BatteryHistory is a simple application that displays an interactive view of a battery discharge curve. Specifically, the battery history of any active RTLS tag can be viewed.

> The following images show the main views of the application.

<table style="border:none;">
  <tr>
    <td><img src="https://user-images.githubusercontent.com/60503970/187513306-b44f0a74-78bf-4862-bd61-2b19c66154e5.png" width="480"/></td>
    <td><img src="https://user-images.githubusercontent.com/60503970/187514732-1eddc0d5-ec95-4fb4-a469-50dbdfe0e73a.png" width="480"/></td>
  </tr>
</table>

## ⚙️ How it works

The application determines the data points for the curve by addressing the XML API of the RTLS controller. The MAC address of a tag is used as a filter. The API returns the entire battery history of the tag after authentication and request. The data is then plotted on a line graph.

## 🎯 Project goals

- [x] Display battery history of any tag
- [x] Provide data in web frontend
- [x] Interactive selection of the period to be viewed
- [ ] Screenshot function for the adjusted display
- [ ] Import all MAC addresses for easy selection

## 📜 Installation guide

### Build from source

From the root of the source tree, run:

```text
go build cmd/main.go
```

### Deploy with Docker

It is recommended to use [docker-compose](https://docs.docker.com/compose/) as it is very convenient. The following example shows a simple deployment without a proxy.

```yaml
version: '3.9'

services:
  battery-history:
    image: plaenkler/battery-history:latest
    container_name: battery-history
    restart: unless-stopped
    ports:
      - 9000:9000
    volumes:
      - ./battery-history:/app/config
```

### Configuration

At first startup, the program creates a config directory relative to the executable file and a `config.yaml` file in it. The first four parameters must be set according to the RTLS server configuration. The `webPort` is the port on which the webserver of BatteryHistory listens.

```yaml
serverAddress: 127.0.0.1
serverPort: "8550"
serverUser: user
serverPassword: password
webPort: "9000"
```


