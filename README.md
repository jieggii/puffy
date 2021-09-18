# puffy
Puffy is an extremely simple unixway GitHub webhook listener for push events which is supposed to be run as systemd service.

## Dependencies
* golang

## Installation
Simply run `install.sh` script to install puffy. The script will:
* Build puffy binary and move it to `/usr/bin/puffy`.
* Create `puffy` directory at `/etc/` and copy its config file to it.
* Copy `puffy.service` to `/etc/systemd/system/` directory.

## Uninstallation
Puffy can be easily uninstalled using `uninstall.sh` script. Just run it.
It will remove all puffy files, including its config and systemd service file 
if you agree in the interactive prompt

## Usage
### Step 1: Configuration
You need to configure puffy at first. Configuration file is in `TOML` format, so, 
primarily get acquainted with [toml specification](https://github.com/kezhuw/toml-spec).

Open `/etc/puffy/config.toml` (this is where it is located by default) with your favourite text editor:
```toml
host = "0.0.0.0"  # optional
port = 8080
endpoint = "/"  # optional

[[repo]]
name = "username/repo-name"
secret = "my super secret"
exec = "/usr/bin/touch /root/hewwo-cutie^^"

[[repo]]
name = "username/repo-name"
secret = "qwerty12345"
exec = "/usr/bin/sh /home/user/scripts/script.sh"
```

Root fields:
| Field      | Type     | Default value | Descripton                                      |
| ---------- | -------- | ------------- | ----------------------------------------------- |
| `host`     | optional | `0.0.0.0`     | Host puffy will listen to                       |
| `port`     | required | (no default)  | Port puffy will listen to                       |
| `endpoint` | optional | `/`           | Endpoint puffy will listen to                   |

Repository fields:
| Field    | Type     | Description                                                             |
| -------- | -------- | ----------------------------------------------------------------------- |
| `name`   | required | Repository name in format `<username>/<reponame>`, e.g: `jieggii/puffy` |
| `secret` | required | Webhook secret which is set in repository settings                      |
| `exec`   | required | Command to be executed when push event is received. It is always recommended to indicate the full path to the binary. E.g: `/usr/bin/bash /home/user/repo/on-push.sh` |


Edit fields and add your repositories.

_Note: you need to restart puffy after every configuration edits._

### Step 2: Running with systemd
When puffy is set up, it's time to use it! Puffy is supposed to be used with systemd, 
but nothing prevents you from running it in the way you want. 
Directly for example, just using the `puffy` command. Anyway, I will show how to use it with systemd.

At first start it:

`sudo systemctl start puffy.service`

You can check its status if you want to make sure everything's fine:

`systemctl status puffy.service`

And after that, you would probably like to `enable` it, so that it starts with the start of the server

`sudo systemctl enable puffy.service`

## TODO
- [ ] Add payload signature validation
