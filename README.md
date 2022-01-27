# Installation

- Install ganache-cli `npm install -g ganache@v7.0.0-alpha.2`
- Install Docker
- Install Go 1.6. Make sure GOPATH, GOROOT are set and go module is turned on.

# Errors you might face.
- This repository currently depends on other private repositories in the sisu-network organization. Run `make configure-git` to tell git to clone those repositories over SSH.

# Run Sisu
## Run with Docker
1. Create a `tmp` folder and copy the `id_rsa` into that folder. Docker will need it to download and build the container.
```
mkdir tmp
cp ~/.ssh/id_rsa tmp
```
2. Run the following command to build all docker images (including dheart and deyes).
```
./scripts/docker_build_all.sh
```
3. Next step is to build sisu and generate config and genesis data. You can replace `1` with any number of node in your network.
```
go build -o ./sisu cmd/sisud/main.go
./sisu local-docker --v 1 --output-dir ./output
```
4. Start docker
```
cd output
docker-compose up -d
```
You can watch logs from sisu, dheart and deyes by running
```
docker-compose logs -f sisu0
```
Replace `0` with any other node index.

## Shutdown Docker
From `output` folder, run:

```
docker-compose down
```

NOTE: if you get into bad data state or you want to reset the blockchain, simply run the local-docker command again from the Sisu root folder
```
./sisu local-docker --v 2 --output-dir ./output
```


## Run without Docker
Running without docker requires 5 different tabs on the terminal. In exchange, you get faster compilation time for Sisu every time you make any change.

## Run ganache-cli
---

Download ganache (make sure you have version **7.x with node 14** and above) and runs the following commands on 2 different terminals:

```
ganache-cli --accounts 10 --blockTime 3 --port 8545 --defaultBalanceEther 100000 --chain.networkId 189985 --chain.chainId 189985 --mnemonic "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum"
```

```
ganache-cli --accounts 10 --blockTime 3 --port 8545 --defaultBalanceEther 100000 --chain.networkId 189986 --chain.chainId 189986 --mnemonic "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum"
```

### Run dheart and deyes
Follow the instruction on dheart and deyes to run these 2 components in 2 separate tabs.

### Build and run Sisu

```
cp .env.dev .env
```

Build Sisu binary file:

```
go build -o ./sisu cmd/sisud/main.go
```

Generate config file and genesis for local sisu app:

```
./scripts/gen_localnet.sh
```

To run the app on localhost:

```
./scripts/run_local.sh
```

You can now deploy ETH transaction on ganache1 at port 7545.

# Interact with Sisu

Sisu has a number of commands for you to interact with the network. At the root folder, type `./sisu` to see a list of command options. Most of them are default commands from Cosmos SDK. We are only interested in `./sisu dev` command for local dev.

After the start of Sisu, wait until the sisu log has the following message similar to the following:

```
adding watcher address 0x1D156a3e1356b58733305e670D61018001997f6E
```

Please note that `0x1D156a3e1356b58733305e670D61018001997f6E` could be replaced by another address. It's the address of the signing key that Sisu network controls. We need to fund that address with some ETH so that Sisu can have balance to deploy contracts and invoke function calls.

Run this command to fund 2 gateway smart contract on 2 different chains ganache1 and ganache2:

```
./sisu dev fund-account ganache1 7545 ganache2 8545 10
```

This commands will fund Sisu's signing key account with 10 ethereum. The last number in the command is the amount you want to fund.

Waits for 10 seconds for the transaction to finalized and for Sisu to deploy its gateway contract.

## Deploy ERC20 token and do token swapping

You can now deploy ERC20 and do token swapping with the repo [smart-contract](https://github.com/sisu-network/smart-contracts).
