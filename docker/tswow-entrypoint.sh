#!/usr/bin/env bash
set -euo pipefail

export NVM_DIR="/root/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh"

if [ ! -d "/tswow-root/tswow-source/.git" ]; then
  echo "Cloning TSWoW sources..."
  git clone https://github.com/tswow/tswow.git --recurse /tswow-root/tswow-source
fi

if [ -f "/tswow-root/tswow-install/TrinityCore/install/trinitycore/bin/authserver" ]; then
  echo "Build artifacts detected at $AUTH_PATH; skipping 'npm run build'."
else
  echo "No build artifacts found; running 'npm run build'..."
  (cd /tswow-root/tswow-source && npm i && npm run build)
fi

cd /tswow-root/tswow-install
tail -f
