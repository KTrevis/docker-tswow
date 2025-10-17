## TSWoW Docker - Quick Start

### 1) First run: provide a WoW 3.3.5 client
On the very first start, you will likely see an error stating that the client cannot be found. You must provide your own World of Warcraft 3.3.5 (Wrath) client:

- Place your client anywhere inside the `tswow-root` folder. For example: `tswow-root/client`.
- Edit `tswow-root/tswow-install/node.conf` and update the client path accordingly. The path is relative to the `tswow-install` directory. For example, if your client is at `tswow-root/client`, then set the path to `../client` in `node.conf`.
- Restart the container after saving `node.conf`.

Once the correct client path is set, the container should start successfully.

### 2) Running commands inside the container
To run commands (e.g., building datascripts/livescripts):

1. Attach to the main process of the container:
   - `docker attach tswow`
2. Type `help` to discover available commands.
3. To detach without stopping the server, press: Ctrl-p then Ctrl-q.

You can also open a separate shell if needed:
- `docker exec -it tswow bash`

### 3) Database credentials and initialization
Database credentials are stored in `.env` and are used when creating the databases for the first time. If you want to change them, do it before the initial creation (i.e., before the first successful start that initializes the DBs). If the DBs already exist and you need to change credentials, you will need to adjust them accordingly in MySQL and in `node.conf`.

### Notes
- The container exposes the default TrinityCore ports: `3724` (authserver) and `8085` (worldserver). Ensure your OS firewall allows inbound TCP for these if you connect from another machine on the LAN.
- You can follow logs with `docker logs -f tswow`.

