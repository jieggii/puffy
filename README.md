# puffy
(WIP) Puffy is an extremely simple unixway GitHub webhook listener and handler for push events

## Todo
- [ ] Add payload signature validation (WIP)

## Dependencies
* golang

## Installation
Simply run `install.sh` script to install puffy. The script will:
* Build puffy binary and move it to `/usr/bin/puffy`
* Create `puffy` directory at `/etc/` and copy its config file to it
* Copy `puffy.service` to `/etc/systemd/system/` directory

## Coming soon
...
