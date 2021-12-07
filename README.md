# puffy
Puffy is an extremely simple GitHub webhook listener for push events which is supposed to be run as systemd service.

## Dependencies
* golang
* make

## Installation
```shell
make install
```
This make command will:
* Build puffy binary and move it to `/usr/bin/puffy`
* Create `puffy` directory at `/etc/` and copy its config file to it
* Copy `puffy.service` to `/etc/systemd/system/` directory

## Uninstallation
```shell
make uninstall
```
All puffy files, including its config directory (`/etc/puffy`) and systemd service file will be removed.

## Usage
### Step 1: Configuring
You need to configure puffy at first. Configuration file is in `TOML` format, so 
primarily get acquainted with [toml specification](https://toml.io/en/v1.0.0) (especially pay attention to [array of tables](https://toml.io/en/v1.0.0#array-of-tables)).

Then open `/etc/puffy/config.toml` (this is where puffy config file is located by default) with your favourite text editor:
```toml
host = "0.0.0.0"  # optional
port = 8080
endpoint = "/"  # optional

[[repo]]
name = "username/repo-name"
exec = "touch file1"

[[repo]]
name = "username/repo-name"
exec = "sh /home/user/scripts/script.sh"

[[repo]]
name = "jieggii/puffy"
exec = "bash /home/jieggii/scripts/alert.sh"
```

Root fields:
| Field      | Type     | Default value | Description                        |
| ---------- | -------- | ------------- | ---------------------------------- |
| `host`     | optional | `0.0.0.0`     | Host to listen to                  |
| `port`     | required | (no default)  | Port to listen to                  |
| `endpoint` | optional | `/`           | Endpoint to listen to              |

Repository fields:
| Field    | Type     | Description                                                                                |
| -------- | -------- | ------------------------------------------------------------------------------------------ |
| `name`   | required | Repository name in format `<username>/<reponame>`, e.g: `jieggii/puffy`                    |
| `exec`   | required | Command to be executed when push event is received, e.g: `bash /home/user/repo/on-push.sh` |

Edit fields and add your repositories.

_Note: you need to restart puffy after every configuration file edits._

### Step 2: Running
When puffy is set up, it's time to run it! Puffy is supposed to be used with systemd, 
but nothing prevents you from running it in the way you want. 
Directly for example, just using the `puffy` command. Anyway, I will show how to use it with systemd.

At first start puffy service:

`sudo systemctl start puffy.service`

You can check its status if you want to make sure everything's fine:

`systemctl status puffy.service`

And after that you would probably like to `enable` it to make puffy start right after boot:

`sudo systemctl enable puffy.service`

## TODO
- [ ] Add payload signature validation
- [ ] Write webhook setup guide
