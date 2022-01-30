all: build

clean:
	rm -rf bin/ vendor/

fmt:
	gofmt -w cmd

build:
	go mod vendor
	mkdir -p bin/
	go build -o puffy cmd/*.go

install: build
	cp -iv puffy /usr/bin/puffy
	mkdir -pv /etc/puffy
	cp -iv config.default.toml /etc/puffy/config.toml
	cp -iv puffy.service /etc/systemd/system/

uninstall:
	@echo '[*] Uninstalling puffy (answer "y" to all prompts to fully uninstall it)...'
	-rm -f /usr/bin/puffy
	-rm -ir /etc/puffy
	-rm -i /etc/systemd/system/puffy.service
