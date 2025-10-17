#!/usr/bin/env bash
set -euo pipefail

# Script to configure database settings in node.conf using environment variables

NODE_CONF_PATH="/tswow-root/tswow-install/node.conf"

# Default values if env vars are not set
DB_USER=${DB_USER:-root}
DB_PASSWORD=${DB_PASSWORD:-password}
DB_HOST=${DB_HOST:-mysql-tswow}
DB_PORT=${DB_PORT:-3306}

# Construct the database connection string
DB_CONNECTION_STRING="${DB_HOST};${DB_PORT};${DB_USER};${DB_PASSWORD}"

echo "Configuring database settings in node.conf..."
echo "Connection string: ${DB_HOST};${DB_PORT};${DB_USER};****"

# Replace Database.WorldSource
sed -i "s|^Database\.WorldSource = .*|Database.WorldSource = \"${DB_CONNECTION_STRING}\"|" "${NODE_CONF_PATH}"

# Replace Database.WorldDest
sed -i "s|^Database\.WorldDest = .*|Database.WorldDest = \"${DB_CONNECTION_STRING}\"|" "${NODE_CONF_PATH}"

# Replace Database.Auth
sed -i "s|^Database\.Auth = .*|Database.Auth = \"${DB_CONNECTION_STRING}\"|" "${NODE_CONF_PATH}"

# Replace Database.Characters
sed -i "s|^Database\.Characters = .*|Database.Characters = \"${DB_CONNECTION_STRING}\"|" "${NODE_CONF_PATH}"

echo "Database configuration complete!"

# After build, ensure world source/dest schemas exist and import latest TDB if present
SQL_FILE="/tswow-root/tswow-build/TDB_full_world_335.24081_2024_08_17.sql"

if [ -f "$SQL_FILE" ]; then
  echo "Preparing databases default.dataset.world.source and default.dataset.world.dest..."

  MYSQL_CMD=(mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}")

  # Create schemas if they don't exist
  "${MYSQL_CMD[@]}" -e "CREATE DATABASE IF NOT EXISTS \`default@002edataset@002eworld@002esource\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
  "${MYSQL_CMD[@]}" -e "CREATE DATABASE IF NOT EXISTS \`default@002edataset@002eworld@002edest\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

  echo "Importing $SQL_FILE into source... (this may take a while)"
  "${MYSQL_CMD[@]}" default@002edataset@002eworld@002esource < "$SQL_FILE"

  echo "Importing $SQL_FILE into dest... (this may take a while)"
  "${MYSQL_CMD[@]}" default@002edataset@002eworld@002edest < "$SQL_FILE"

  echo "SQL import complete."
else
  echo "SQL file not found: $SQL_FILE (skipping import)"
fi