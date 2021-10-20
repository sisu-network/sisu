## Installation

Install Go 1.6. Make sure GOPATH, GOROOT are set and go module is turned on.
A few depdencies of this project are private repo. You need to config git to use `git` instead of `https` when downloading repo.
Run the following command to replace `https` by `git`.

```
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

To confirm, do `more ~/.gitconfig` and make sure you see the following:

```
[url "git@github.com:"]
	insteadOf = https://github.com/
```

## Running locally without TSS

From the sisu root folder

```
cp .env.dev .env
```

Generate config file and genesis for local sisu app:

```
./scripts/gen_testnet.sh
```

Disable the TSS component in the file `~/.sisu/main/config/sisu.toml` by changing the default tss settings to false.

```
[tss]
enable = false
```

To run the app on localhost:

```
./scripts/run_local.sh
```

You can now deploy ETH transaction on Sisu at port 1234.

## Running locally wihth TSS

You need to enable tss in `~/.sisu/main/config/sisu.toml` in order to run full Sisu node.

### Running without docker

You will need at least 5 different terminal tabs to run a full Sisu nodes: 2 tabs for ganache-cli, 1 for dheart, 1 for deyes, 1 for sisu.

#### Run ganache-cli

Download ganache-cli (make sure you have version 6.x) and runs the following commands on 2 different terminals:

```
ganache-cli --accounts 10 --blockTime 3 --port 7545 --defaultBalanceEther 100000 --chainId 1 --mnemonic "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum"
```

```
ganache-cli --accounts 10 --blockTime 3 --port 8545 --defaultBalanceEther 100000 --chainId 36767 --mnemonic "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum"
```

These commands create 2 simulated blockchains on port 7545 and 8545.

#### Run dheart and deyes

Runs dheart and deyes using their instruction in the repos.

#### Run Sisu

```
./scripts/run_local.sh
```
