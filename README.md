# puffy
**puffy** is an extremely simple GitHub webhook listener for **push events** 
which is supposed to be run as **systemd service**.

## Dependencies
* **golang**
* **sh** (used by default) or any other shell that treats `-c` flag as command to execute
* **make** - optinal, used only for installation  

## Installation
The program can be easily installed using **make** command:

```shell
make install
```

It will:
* Build **puffy** binary and move it to `/usr/bin/puffy`
* Create `/etc/puffy/` directory and copy [example config file](https://github.com/jieggii/puffy/blob/master/config.example.toml) to it
* Copy `puffy.service` to `/etc/systemd/system/` directory

## Uninstallation
To uninstall **puffy** simply run:

```shell
make uninstall
```

All **puffy** data (its binary - `/usr/bin/puffy`, config directory - `/etc/puffy/` 
and **systemd** service file - `/etc/systemd/system/puffy.service`) will be removed. 
All deletions will require your confirmation.

## Usage guide
<details>
<summary>Step 1: configure puffy</summary>

You will need to configure **puffy** at first. Configuration file is in **TOML** format, so 
primarily get acquainted with [toml specification](https://toml.io/en/v1.0.0) 
(especially pay attention to [array of tables](https://toml.io/en/v1.0.0#array-of-tables)).

Then open `/etc/puffy/config.toml` (this is where its config file is located by default)
with your favourite text editor:
```toml
host = "0.0.0.0"  # host to listen to (default: "0.0.0.0")
port = 8080       # port to listen to
endpoint = "/"    # endpoint to listen to (default: "/")

# shell to use when running command from $repo.exec 
shell = "/usr/bin/bash"  # default: "/usr/bin/sh"

# directory to go to before executing command from $repo.exec
workdir = "/"  # default: "/"

[[repos]]  # full repository example
name = "username/repo-name"            # name of the repository in <username>/<repo-name> format
shell = "/usr/bin/fish"                # (optional) overwrites $shell for this repository
workdir = "/home/username/repo-name/"  # (optional) overwrites $workdir for this repository
exec = "./script.fish"                 # command to execute when push event is received

[[repos]]  # the most simple repository example
name = "username/repo-name"
exec = "/home/username/scripts/alert.sh"

# other repository examples
[[repos]]  
name = "username/repo-name"
workdir = "/home/username/repo-name/"
exec = "git pull"

[[repos]]
name = "username/website"
workdir = "/home/username/repos/website/"
exec = "bash scripts/on-push.bash"
```

Edit fields and add your repositories.

_**Note:** you need to restart puffy after every config file edits._
</details>

<details>
<summary>Step 2: run puffy</summary>

When **puffy** is set up, it's time to run it! **Puffy** is supposed to be used with **systemd**, 
but nothing prevents you from running it in the way you want. 
Directly for example, just using the `puffy` command. And I recommend you to do it at first just to make sure everything's fine. Anyway, I will show how to use it with **systemd**.

At first start the puffy service:

`sudo systemctl start puffy.service`

You can check its status if you want to make sure it is running properly:

`systemctl status puffy.service`

And after that you would probably like to *enable* it to make puffy always start after boot:

`sudo systemctl enable puffy.service`

You can also read puffy logs using this command:

`sudo journalctl -u puffy.service`
</details>

<details>
<summary>Step 3: set up your GitHub repository</summary>

Now, when **puffy** is configured, running on your servier and waiting for push events,
it's time to configure your GitHub repository.

1. Go to repository **settings** and choose **Webhooks** meny entry.
![pic1](https://imgur.com/To3W0yT.jpg)

2. Press **Add webhook** and confirm your password.
3. Fill fields: Provide **payload URL** in `http://<hostname>:<port>/<endpoint>` format, where `<hostname>` is your domain name or IP address, `<port>` and `<endpoint>` are port and endpooint **puffy** is listening to; set **Content type** to `application/json` and press **Add webhook**.
![pic2](https://imgur.com/tKDBryR.jpg)

**Done!** Webhook is configured. Now, to check if everythng works fine, 
click on your webhook, then go to **Recent deliveries** tab and click on the first delivery. 
It should look like this (with response code **200** and `pong!` body):
![pic3](https://imgur.com/inL7aXG.jpg)
</details>

## TODO
- [ ] Add payload signature validation
