#!/usr/bin/env bash
echo "Uninstalling puffy..."
sudo rm -f /usr/bin/puffy
sudo rm -ir /etc/puffy
sudo rm -i /etc/systemd/system/puffy.service
echo "Done!"