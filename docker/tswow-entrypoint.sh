#!/usr/bin/env bash
set -euo pipefail

export NVM_DIR="/root/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh"

git config --global safe.directory '*'

if [ ! -d "/tswow-root/tswow-source/.git" ]; then
  echo "Cloning TSWoW sources..."
  git clone https://github.com/tswow/tswow.git --recurse /tswow-root/tswow-source
fi

# Inject a cmake wrapper to disable libpng NEON globally without modifying sources
cat >/usr/local/bin/cmake <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
exec /usr/bin/cmake \
  -DCMAKE_C_FLAGS="-DPNG_ARM_NEON=0 -DPNG_ARM_NEON_OPT=0 ${CMAKE_C_FLAGS:-}" \
  -DCMAKE_CXX_FLAGS="-DPNG_ARM_NEON=0 -DPNG_ARM_NEON_OPT=0 ${CMAKE_CXX_FLAGS:-}" \
  "$@"
EOF
chmod +x /usr/local/bin/cmake

# Disable libpng NEON to avoid undefined reference to png_init_filter_functions_neon on non-ARM builds
export CFLAGS="${CFLAGS:-} -DPNG_ARM_NEON=0"
export CXXFLAGS="${CXXFLAGS:-} -DPNG_ARM_NEON=0"
export CPPFLAGS="${CPPFLAGS:-} -DPNG_ARM_NEON=0"
export CMAKE_C_FLAGS="${CMAKE_C_FLAGS:-} -DPNG_ARM_NEON=0"
export CMAKE_CXX_FLAGS="${CMAKE_CXX_FLAGS:-} -DPNG_ARM_NEON=0"

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
