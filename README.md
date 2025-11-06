<pre>
<code>
______ _     _   ___   __ _____ ___________ _____  ___  ___  ___
|  ___| |   | | | \ \ / //  ___|_   _| ___ \  ___|/ _ \ |  \/  |
| |_  | |   | | | |\ V / \ `--.  | | | |_/ / |__ / /_\ \| .  . |
|  _| | |   | | | |/   \  `--. \ | | |    /|  __||  _  || |\/| |
| |   | |___| |_| / /^\ \/\__/ / | | | |\ \| |___| | | || |  | |
\_|   \_____/\___/\/   \/\____/  \_/ \_| \_\____/\_| |_/\_|  |_/
</code>
</pre>
</br>
<center> Fluxstream - Open-Source Torrent Streamer.
</center>
<br>
FluxStream is an open-source, self-hosted streaming platform built around torrent-based media delivery.
It lets you stream content instantly from magnet links.

---
### Quick Install

**Prerequisite:**  
> Docker **must** be installed and running on your system.  
>
> [![Install Docker](https://img.shields.io/badge/Install%20Docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)](https://docs.docker.com/get-started/get-docker/)

##### Linux/macOS
> Copy and paste this in your terminal:
```sh
$ curl -fsSL https://raw.githubusercontent.com/scythe504/fluxstream/main/scripts/install.sh | bash
```
##### Windows
> Copy and paste this in your terminal:
```powershell
irm https://raw.githubusercontent.com/scythe504/fluxstream/main/scripts/install.ps1 | iex
```

#### Verify Installation (Only for Quick Install)

Once installation completes, verify that fluxstream was installed correctly:

```sh
$ fluxstream --version
```
It should print something like this in your terminal.
```sh
fluxstream version 0.1.3
```
To view all available commands
```sh
$ fluxstream 
```
or
```sh
$ fluxstream --help 
```
#### Run Fluxstream (Only for Quick Install)
This will generate docker-compose.yml and download directories for the media.
```sh
$ fluxstream setup
```
This will spin up the docker containers for backend(Server) and the frontend (Web Client).
```sh
$ fluxstream start
```
Stopping the containers.
```sh
$ fluxstream stop
```
or
```sh
$ docker ps
$ docker kill <CONTAINER-ID>
```

Prints where the web client is running. 
```sh
$ fluxstream where
```
Expected Output:
```sh
#  --- May happen on Windows (using WSL) ---
FluxStream web interface available at:
  Local: http://localhost:3000
  Network: http://172.17.205.43:3000

 WARN: If Network IP is 172.*, it will not open in other devices.
```

---
### Manual Install

**Prerequisites:**
> Requires Go and Node.js installed on your system.  
>
> [![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/dl/) [![Node.js](https://img.shields.io/badge/Node.js-%23339933.svg?style=for-the-badge&logo=node.js&logoColor=white)](https://nodejs.org/en/download/)

##### Backend (Server)

```sh
$ git clone https://github.com/scythe504/fluxstream.git
$ cd fluxstream
$ cp .env.example .env
```
```sh
# For cmd/api/main.go only
$ make run 
```
or
```sh
$ go run cmd/api/main.go
```
The Server will be up on 
```
http://localhost:8080
```

##### Frontend (Web Client)

```sh
$ git clone https://github.com/scythe504/fluxstream-web.git
$ cd fluxstream-web
$ cp .env.example .env
```

```sh
$ npm install
$ npm run build
$ npm run start
```
---

### License

FluxStream is open source software licensed under the  
**GNU Affero General Public License v3.0 (AGPLv3)**.  
See [LICENSE](./LICENSE) for details.