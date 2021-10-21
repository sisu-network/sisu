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
./scripts/gen_localnet.sh
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

In addition, you need to install mysql locally with the following config:

```
host: localhost
port: 3306
username: root
password: password
```

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

#### Interact with Sisu

Sisu has a number of commands for you to interact with the network. At the root folder, type `./sisu` to see a list of command options. Most of them are default commands from Cosmos SDK. We are only interested in `./sisu dev` command for local dev.

After the start of Sisu, wait until the sisu log has the following message similar to the following:

```
adding watcher address 0x1D156a3e1356b58733305e670D61018001997f6E
```

Please note that `0x1D156a3e1356b58733305e670D61018001997f6E` could be replaced by another address. It's the address of the signing key that Sisu network controls. We need to fund that address with some ETH so that Sisu can have balance to deploy contracts and invoke function calls.

```
./sisu dev fund-account eth 7545 sisu-eth 8545 10
```

This commands will fund Sisu's signing key account with 10 ethereum. The last number in the command is the amount you want to fund.

Waits for 10 seconds for the transaction to finalized and for Sisu to deploy its gateway contract.

Now you can start transferring ERC20 tokens out of a blockchain. You need a deployed ERC20 contract on one of the ganache chain. You can deploy using separate service or use Sisu command line:

```
./sisu dev deploy-erc20 eth
```

Waits for a few seconds and you will see the address the newly deploy ERC20 contract. You can transfer token to second dev chain using this command format:

```
./sisu dev transfer-out [ContractType] [FromChain] [TokenAddress] [ToChain] [RecipientAddress]
```

For example

```
./sisu dev transfer-out erc20 eth 0xB369Be7F62cfb3F44965db83404997Fa6EC9Dd58 sisu-eth 0xE8382821BD8a0F9380D88e2c5c33bc89Df17E466
```

Waits a few seconds for the transaction to complete. Afterward, you can query the asset balance on the destination chain:

```
./sisu dev query [ContractType] [chain] [AssetId] [AccountAddress]

./sisu dev query erc20 sisu-eth eth__0xB369Be7F62cfb3F44965db83404997Fa6EC9Dd58 0xE8382821BD8a0F9380D88e2c5c33bc89Df17E466
```
