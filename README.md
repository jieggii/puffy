# (WIP) puffy
Puffy is an extremely simple unixway GitHub webhook listener and handler for push events

## Todo
- [ ] Add payload signature validation (WIP)

## Dependencies
* golang

## Installation
Simply run `install.sh` script to install puffy. The script will:
* Build puffy binary and move it to `/usr/bin/puffy`
* Create `puffy` directory at `/etc/` and copy its config file to it
* Copy `puffy.service` to `/etc/systemd/system/` directory

## Running
### 1: Configuration
First you need to configure puffy. Configuration file is in `TOML`, so, 
primarily get acquainted with [toml specification](https://github.com/kezhuw/toml-spec)

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
exec = "/home/user/scripts/script.sh"
```

Edit fields as you want! A new repository can be added under `[[repo]]` line.
Please be sure to indicate `name`, `secret` and `exec` fields.

_Also note, that after every configuration edits you need to restart puffy_

### 2: Running with systemd
When everything is set up, it's time to use the program what it is for was made.
Puffy is in general meant to be used as a systemd service.
At first start it:

`sudo systemctl start puffy.service`

You can check its status if you want to make sure everything's fine:

`systemctl status puffy.service`

And after that, you probably would like to `enable` it, so that it starts with the start of the server

`sudo systemctl enable puffy.service`

...