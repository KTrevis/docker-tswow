## TSWoW Docker - Quick Start

### Prerequisites
- Docker and Docker Compose installed
- A World of Warcraft 3.3.5 (Wrath of the Lich King) client

### Setup Instructions

#### 1) Configure database credentials (optional)
Database credentials are stored in `.env` and are used when creating the databases for the first time. If you want to change them, **do it before the initial setup** (i.e., before running `docker compose up` for the first time).

#### 2) Start the containers for the first time
```bash
docker compose up -d
```

The container will build TSWoW and create the necessary directory structure. You can follow the logs:
```bash
docker logs -f tswow
```

**Expected behavior:** On the first run, you will see an error stating that the WoW client cannot be found. This is normal. Continue to step 3.

#### 3) Provide your WoW 3.3.5 client
Once the container has created the directory structure:

1. Place your WoW 3.3.5 client anywhere inside the `tswow-root` folder. For example: `tswow-root/client`.
2. Edit `tswow-root/tswow-install/node.conf` and locate the client path setting.
3. Update the path to point to your client. **The path is relative to the `tswow-install` directory.**
   - Example: if your client is at `tswow-root/client`, set the path to `../client` in `node.conf`.
4. Restart the container:
   ```bash
   docker compose restart
   ```

Once the correct client path is set, the container should start successfully and build the server.

#### 4) Running commands inside the container
To run TSWoW commands (e.g., building datascripts/livescripts):

1. Attach to the main process of the container:
   ```bash
   docker attach tswow
   ```
2. Type `help` to discover available commands.
3. To detach without stopping the server, press: **Ctrl-p** then **Ctrl-q**.

