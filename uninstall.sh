#!/usr/bin/env bash

echo "Removing FletaloYa scripts from '/usr/local/bin'"

for f in $(find scripts -type f -execdir echo '{}' ';')
do
    echo Removing /usr/local/bin/$f
    rm /usr/local/bin/$f
done

echo "For installing again just run './install.sh'"
