#!/usr/bin/env bash
set -euo pipefail

export NVM_DIR="/opt/nvm"
[ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh"

git config --global safe.directory '*'

if [ ! -d "/tswow-root/tswow-source" ]; then
  echo "Cloning TSWoW sources..."
  git clone https://github.com/tswow/tswow.git --recurse /tswow-root/tswow-source
fi

if [ -f "/tswow-root/tswow-install/bin/trinitycore/RelWithDebInfo/authserver" ]; then
  echo "Build artifacts detected; skipping 'npm run build'."
else
  echo "No build artifacts found; running 'npm run build'..."
  (cd /tswow-root/tswow-source && npm i && npm run build)
fi

cd /tswow-root/tswow-install
# TODO: remplacer ce sleep de golmon par un truc qui attend que la db soit up
sleep 10

/docker/tswow-build.sh

npm run start
