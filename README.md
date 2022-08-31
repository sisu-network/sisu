# Installation

- Install ganache-cli `npm install -g ganache@v7.0.0-alpha.2`
- Install Docker
- Install Go 1.6. Make sure GOPATH, GOROOT are set and go module is turned on.

# Run Sisu

## Run without Docker
Running without docker requires 5 different tabs on the terminal. In exchange, you get faster compilation time for Sisu every time you make any change.

## Run ganache-cli
---

Download ganache (make sure you have version **7.x with node 14** and above) and runs the following commands on 2 different terminals:

```
ganache-cli --accounts 10 --blockTime 3 --port 7545 --defaultBalanceEther 100000 --chain.networkId 189985 --chain.chainId 189985 --mnemonic "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum"
```

```
ganache-cli --accounts 10 --blockTime 3 --port 8545 --defaultBalanceEther 100000 --chain.networkId 189986 --chain.chainId 189986 --mnemonic "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum"
```

### Run dheart and deyes
Follow the instruction on dheart and deyes to run these 2 components in 2 separate tabs. Make sure you create `dheart.toml` and `deyes.toml` files before running them.

### Deploy Vault contract
You need to deploy Vault contract before running Sisu. This vault is later controlled by Sisu's network account .

Build Sisu binary file:

```
go build -o ./sisu cmd/sisud/main.go
```

Deploy vault contract:
```
./sisu dev deploy
```

### Build and run Sisu

Generate config file and genesis for local sisu app. You need to generate this while ganache instances are running.

```
./scripts/gen_localnet.sh
```

Create `.env` file

```
cp .env.dev .env
```

Build Sisu binary file:

```
go build -o ./sisu cmd/sisud/main.go
```

To run the app on localhost:

```
./scripts/run_local.sh
```

You can now deploy ETH transaction on ganache1 at port 7545.

# Interact with Sisu

Run this command to fund 2 gateway smart contract on 2 different chains ganache1 and ganache2:

```
./sisu dev deploy-and-fund
```

This command does a number of things:
- Deploy ERC20 contracts in 2 ganache blockchains.
- Fund the Sisu's network public key with some ETH (the last argument from command line)
- Wait for Sisu to finish deploying gateway contracts
- Transfer ERC20 to gateway contracts

After this command finishes running, you can start swapping token through both chains.

```
./sisu dev swap --erc20-symbol SISU --amount 10 --account 0x2d532C099CA476780c7703610D807948ae47856A
```

Wait for few seconds for Sisu to pick up and execute the swap. Afterward, you can query to verify that the recipient does have some ERC20 token in the destination chain (ganache2).

```
./sisu dev query --erc20-symbol SISU --chain ganache2 --account 0x2d532C099CA476780c7703610D807948ae47856A
```

Note: the token the recipient receives will not exactly the amount user swaps because there is some fee taken away from the swap amount.


## Run with Docker
1. Run the following command to build all docker images (including dheart and deyes).
```
./scripts/docker_build_all.sh
```
2. Next step is to build sisu and generate config and genesis data. You can replace `2` with any number of node in your network.
```
go build -o ./sisu cmd/sisud/main.go
./sisu local-docker --v 2 --output-dir ./output
```
3. Start docker
```
cd output
docker-compose up -d
```
You can watch logs from sisu, dheart and deyes by running
```
docker-compose logs -f sisu0
```
You can view log of other nodes by replacing `0` with any other node index.

4. Deploy vault contract

```
./sisu dev deploy-and-fund
```

5. Deploy ERC20 & fund vault with some tokens:

```
./sisu dev deploy-and-fund
```

6. Swapping SISU token between ganache1 & ganache2:

```
./sisu dev swap --erc20-symbol SISU --amount 10 --account 0x2d532C099CA476780c7703610D807948ae47856A
```

7. Query recipient's balance in the destination chain:

```
./sisu dev query --erc20-symbol SISU --chain ganache2 --account 0x2d532C099CA476780c7703610D807948ae47856A
```


## Shutdown Docker
From `output` folder, run:

```
docker-compose down
```

NOTE: if you get into bad data state or you want to reset the blockchain, simply run the local-docker command again from the Sisu root folder
```
./sisu local-docker --v 2 --output-dir ./output
```
