all: build

clean:
	@echo "[*] Cleaning directory..."
	rm -rf bin/ vendor/

fmt:
	gofmt -w cmd

build:
	@echo "[*] Building puffy..."
	go mod vendor
	mkdir -p bin/
	go build -o bin/puffy cmd/main.go
	@echo "[+] Puffy was built. Binary: ./bin/puffy"

install: build
	@echo "[*] Installing puffy..."
	@sudo mv -iv ./bin/puffy /usr/bin/puffy
	@sudo mkdir -pv /etc/puffy
	@sudo cp -iv config.example.toml /etc/puffy/config.toml
	@sudo cp -iv puffy.service /etc/systemd/system/
	@echo "[+] Puffy was successfully installed"
	@echo "Now you would probably like to edit the config file (/etc/puffy/config.toml)"
	@echo "For more information please see README.md at https://github.com/jieggii/puffy"

uninstall:
	@echo '[*] Uninstalling puffy (answer "y" to all prompts to fully uninstall it)...'
	@-sudo rm -f /usr/bin/puffy
	@-sudo rm -ir /etc/puffy
	@-sudo rm -i /etc/systemd/system/puffy.service
	@echo "[+] Puffy was uninstalled"

help:
	@echo "Command     Description      "
	@echo "-----------------------------"
	@echo "build       build puffy"
	@echo "install     install puffy"
	@echo "uninstall   uninstall puffy"
	@echo "fmt         format sources"
	@echo "clean       clean project dir"
