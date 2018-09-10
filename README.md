# Dr. Vault | Inspired by Dr. Fate | [Docker Hub](https://hub.docker.com/r/reeganexe/dr-vault/)
<center>
  <img src="https://pm1.narvii.com/6767/f75f79b540a7648a7e618311bb12d1d2b31c41b3v2_00.jpg" alt="Dr. Fate - DC Comic" />
</center>

Monitoring a directory and sync all YAML files to Vault.

Use for development purpose.

```sh
go get github.com/reeganexe/dr-vault
```

### Usage
```sh
dr-vault --vault-address '0.0.0.0:8200' --vault-token root --dir $PWD/your-configs-dir
```

### Help
```
NAME:
   Dr. Vault - Vault folder monitoring

USAGE:
   dr-vault [options...]

EXAMPLE:
   dr-vault --vault-address '0.0.0.0:8200' --vault-token root --dir $PWD/your-configs-dir

AUTHOR:
   Ninh Pham #ReeganExE -> ninh.js.org

OPTIONS:
   --vault-address value, -a value  Vault address that dr-vault will connect to (default: "0.0.0.0:8200") [$VAULT_DEV_LISTEN_ADDRESS]
   --vault-token value, -t value    A writable token (default: "root") [$VAULT_DEV_ROOT_TOKEN_ID]
   --dir value, -d value            Specify a directory to monitor. (default: "/var/source") [$MONITOR_DIR]
   --verbose, -p                     [$VERBOSE]
   --help, -h                       show help
   --version, -v                    print the version
```

## Docker

I already built a docker image that included Vault Dev Server and Dr. Vault

Define vault address (`VAULT_DEV_LISTEN_ADDRESS`) and root token (`VAULT_DEV_ROOT_TOKEN_ID`) then mount your source directory to `/var/source`.

```sh
docker run --rm -p 8200:8200 \
  -e 'VAULT_DEV_ROOT_TOKEN_ID=root' \
  -e 'VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8200' \
  -v '/path/to/your/dir:/var/source' \
  --name "configuration-service.infra" \
  reeganexe/dr-vault:0.1
```

### Sample docker-compose.yml

```yml
version: "3"
services:
  vault:
    image: reeganexe/dr-vault:0.1
    environment:
      - "VAULT_DEV_ROOT_TOKEN_ID=root"
      - "VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8200"
    ports:
      - "8200:8200"
    volumes:
      - $PWD/test-dir:/var/source
```

### MONITOR_DIR
By default, `dr-vault` will watch on the directory `/var/source` (inside the container), you can change to any others via `MONITOR_DIR` environment.

```sh
docker run --rm -p 8200:8200 \
  -e 'VAULT_DEV_ROOT_TOKEN_ID=root' \
  -e 'VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8200' \
  -e 'MONITOR_DIR=/your/path' \
  -v '/test-dir:/your/path' \
  reeganexe/dr-vault:0.1
```
