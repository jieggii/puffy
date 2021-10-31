build:
	@echo "[*] Building puffy..."
	go mod download github.com/BurntSushi/toml
	go build src/main.go
	@echo "Puffy binary: ./main"

install: build
	@echo "[*] Installing puffy..."
	@sudo mv -iv main /usr/bin/puffy
	@sudo mkdir -pv /etc/puffy
	@sudo cp -iv config.example.toml /etc/puffy/config.toml
	@sudo cp -iv puffy.service /etc/systemd/system/
	@echo "[+] Puffy was successfully installed!"
	@echo "Now you would probably like to edit the config file 2(/etc/puffy/config.toml)"
	@echo "For more information please check https://github.com/jieggii/puffy/blob/master/README.md"

uninstall:
	@echo "[*] Uninstalling puffy (answer \"y\" to all prompts to fully uninstall it)"
	@sudo rm -f /usr/bin/puffy
	@sudo rm -ir /etc/puffy
	@sudo rm -i /etc/systemd/system/puffy.service
	@echo "[*] Puffy was successfully uninstalled!"

help:
	@echo "build - build puffy"
	@echo "install - install puffy"
	@echo "uninstall - uninstall puffy"
