Running a full Sisu node is a fairly complicated process even on localhost. To run on localhost, you will need 4 different terminal tabs. There are running terminals for sisu, dheart, deyes and ganache-cli.

# Pre-run steps

### Creates .env configs files.

Create .env config files for each of the repo: sisu, dhearts, deyes and put them in the root folders. The followings are config files for each of the repo. Please replace directories in config files with your directories.

Sisu

```
AES_KEY_HEX=c787ef22ade5afc8a5e22041c17869e7e4714190d88ecec0a84e241c9431add0
```

dhearts

```
USE_ON_MEMORY=true
HOME_DIR=~/.sisu/dheart
DHEART_PORT=28300
BOOTSTRAP_PEERS=192.168.0.2
SISU_SERVER_URL=http://localhost:25456
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=password
DB_SCHEMA=dheart
DB_MIGRATION_PATH=file://db/migrations/
AES_KEY_HEX=c787ef22ade5afc8a5e22041c17869e7e4714190d88ecec0a84e241c9431add0
```

deyes

```
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=password
DB_SCHEMA=deyes
CHAIN_RPC_URL=http://localhost:7545
SERVER_PORT=31001
SISU_SERVER_URL=http://localhost:25456
BLOCK_TIME=1000
STARTING_BLOCK=0
CHAIN=eth
```

### Generate data for Sisu localnet

Create `.sisu` folder in your home directory

```
mkdir -p ~/.sisu
```

In sisu repo root folder, runs:

```
./scripts/gen_localnet.sh
```

### Enable TSS for sisu

You need to enable TSS feature for full Sisu node. To do this, open the `~/.sisu/tss/tss.toml` config file and copy paste the following content:

```
enable = true
host = "localhost"
port = 5678
[supported_chains]
[supported_chains.eth]
  symbol = "eth"
  id = 1
  url = "http://localhost:31001"
```

# Run all programs

You will need a docker-compose file or 4 different terminal tabs to run a full node.

1. Run [ganache-cli](https://www.npmjs.com/package/ganache-cli):

```
ganache-cli --accounts 10 --blockTime 3 --port 7545 --defaultBalanceEther 100000 --mnemonic "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum"
```

2. Run dheart:
   At the root project of dheart repo, simply run the main.go file:

```
go run main.go
```

3. Run deyes:
   At the root project of deyes repo, simply run the main.go file:
   ```
   go run main.go
   ```
4. Run Sisu
   At the root of Sisu project, run:

```
./scripts/run_local.sh
```
