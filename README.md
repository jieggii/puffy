# puffy
**puffy** is an extremely simple GitHub webhook listener for **push events** 
which is supposed to be run as **systemd service**.

## Dependencies
* [golang](https://go.dev)
* [make](https://www.gnu.org/software/make/)

## Installation
The program can be easily installed using **make** command:

```shell
make install
```

It will:
* Build **puffy** binary and move it to `/usr/bin/puffy`
* Create `/etc/puffy/` directory and copy [example (default) **puffy** config file]() to it
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
Contents:
* [Step 1: configuring **puffy**]()
* [Step 2: running **puffy**]()
* [Step3: seting-up **GitHub repository**]()

### Step 1: configuring puffy
You will need to configure **puffy** at first. Configuration file is in **TOML** format, so 
primarily get acquainted with [toml specification](https://toml.io/en/v1.0.0) 
(especially pay attention to [array of tables](https://toml.io/en/v1.0.0#array-of-tables)).

Then open `/etc/puffy/config.toml` (this is where **puffy** config file is located by default) 
with your favourite text editor:
```toml

```

Edit fields and add your repositories.

_Note: you need to restart **puffy** after every config file edits._

### Step 2: running puffy
When **puffy** is set up, it's time to run it! **Puffy** is supposed to be used with **systemd**, 
but nothing prevents you from running it in the way you want. 
Directly for example, just using the `puffy` command. And I recommend you to do it at first just to make sure everything's fine. Anyway, I will show how to use it with **systemd**.



At first start **puffy** service:

`sudo systemctl start puffy.service`

You can check its status if you want to make sure it is running properly:

`systemctl status puffy.service`

And after that you would probably like to *enable* it to make **puffy** always start after boot:

`sudo systemctl enable puffy.service`

### Step 3: setting up your GitHub repository
Now, when **puffy** is configured, running on your servier and waiting for push events it's time to configure your GitHub repository.

## TODO
- [ ] Add payload signature validation
