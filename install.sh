echo "Installing puffy!"
echo "Building puffy..."
go build src/main.go
echo "Copying some files..."
sudo mv -iv main /usr/bin/puffy
sudo mkdir -pv /etc/puffy
sudo cp -iv config.example.toml /etc/puffy/config.toml
sudo cp -iv puffy.service /etc/systemd/system/
echo "Done!"
echo "If you don't see any errors above, puffy is likely installed :)"
echo "Now I suppose you would like to edit the config file (by default it is located at /etc/puffy/config.toml)"
echo "For more information please check https://github.com/jieggii/puffy/README.md"