echo "Uninstalling puffy..."
sudo rm -f /usr/bin/puffy
sudo rm -rf /etc/puffy
sudo rm -f puffy.service /etc/systemd/system/
echo "Done!"