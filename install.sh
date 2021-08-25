echo "Installing puffy..."
go build src/main.go
sudo mv -iv main /usr/bin/puffy
sudo mkdir -pv /etc/puffy
sudo cp -iv config.example.toml /etc/puffy/config.toml
sudo cp -iv puffy.service /etc/systemd/system/
echo "Done!"