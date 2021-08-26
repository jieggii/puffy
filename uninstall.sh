#!/usr/bin/env bash
echo "Uninstalling puffy..."
sudo rm -f /usr/bin/puffy
sudo rm -rif /etc/puffy
sudo rm -if /etc/systemd/system/puffy.service
echo "Done!"