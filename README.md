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

## Running locally

From the sisu root folder

```
cd src
echo "MODE=dev" > .env
```

Generate config file and genesis for local sisu app:

```
./scripts/gen_testnet.sh
```

To run the app on localhost:

```
./scripts/run_local.sh
```
