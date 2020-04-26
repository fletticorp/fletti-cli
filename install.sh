#!/usr/bin/env bash

echo "Installing FletaloYa scripts into '/usr/local/bin'"

for f in $(find scripts -type f -execdir echo '{}' ';')
do
    echo Copying scripts/$f to /usr/local/bin/$f
    cp scripts/$f /usr/local/bin/$f
done

echo "For uninstalling just run './uninstall.sh'"
