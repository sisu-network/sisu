package dev

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"path/filepath"
	"time"

	etypes "github.com/ethereum/go-ethereum/core/types"
	"google.golang.org/grpc"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	econfig "github.com/sisu-network/deyes/config"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/contracts/eth/erc20"
	"github.com/sisu-network/sisu/utils"
	hdwallet "github.com/sisu-network/sisu/utils/hdwallet"
	tssTypes "github.com/sisu-network/sisu/x/sisu/types"
)

const (
	defaultMnemonic = "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum"
	Blocktime       = time.Second * 3
)

var (
	localWallet *hdwallet.Wallet
	account0    accounts.Account
	privateKey0 *ecdsa.PrivateKey
	nonceMap    map[string]*big.Int
)

func init() {
	var err error
	localWallet, err = hdwallet.NewFromMnemonic(defaultMnemonic)
	if err != nil {
		panic(err)
	}

	path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", 0))
	account0, err = localWallet.Derive(path, true)
	if err != nil {
		panic(err)
	}

	privateKey0, err = localWallet.PrivateKey(account0)
	if err != nil {
		panic(err)
	}

	nonceMap = make(map[string]*big.Int)
}

func getPrivateKey(mnemonic string) (*ecdsa.PrivateKey, common.Address) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		panic(err)
	}

	dpath, err := accounts.ParseDerivationPath("m/44'/60'/0'/0/0")
	if err != nil {
		panic(err)
	}

	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)

	key := masterKey
	for _, n := range dpath {
		key, err = key.Derive(n)
		if err != nil {
			panic(err)
		}
	}

	privateKey, err := key.ECPrivKey()
	if err != nil {
		panic(err)
	}

	privateKeyECDSA := privateKey.ToECDSA()
	publicKey := privateKeyECDSA.PublicKey
	addr := crypto.PubkeyToAddress(publicKey)

	return privateKeyECDSA, addr
}

func getSigner(client *ethclient.Client) etypes.Signer {
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	return etypes.NewLondonSigner(chainId)
}

func getEthClients(chains []string, genesisFolder string) []*ethclient.Client {
	clients := make([]*ethclient.Client, 0)
	if len(chains) == 0 {
		return clients
	}

	deyesCfg := econfig.Load(filepath.Join(genesisFolder, "deyes.toml"))
	deyesChains := deyesCfg.Chains

	for _, chain := range chains {
		found := false
		for _, c := range deyesChains {
			if c.Chain == chain {
				for _, rpcUrl := range c.Rpcs {
					client, err := ethclient.Dial(rpcUrl)
					if err != nil {
						log.Errorf("Failed to dial %s on chain %s", rpcUrl, chain)
					} else {
						// Do a sanity call to make sure that this url actually works.
						num, err := client.BlockNumber(context.Background())
						if err == nil && num > 0 {
							log.Verbosef("Use url %s for chain %s", rpcUrl, chain)
							clients = append(clients, client)
							found = true
							break
						} else {
							log.Error("Cannot call to blocknumber, err = ", err)
						}
					}
				}
			}
		}

		if !found {
			panic(fmt.Errorf("Cannot find healthy url for chain %s", chain))
		}
	}

	return clients
}

func queryErc20Balance(client *ethclient.Client, tokenAddr string, target string) (*big.Int, error) {
	store, err := erc20.NewErc20(common.HexToAddress(tokenAddr), client)
	if err != nil {
		return nil, err
	}

	balance, err := store.BalanceOf(nil, common.HexToAddress(target))

	return balance, err
}

func approveAddress(client *ethclient.Client, mnemonic string, erc20Addr string, target string) {
	contract, err := erc20.NewErc20(common.HexToAddress(erc20Addr), client)
	if err != nil {
		panic(err)
	}

	opts, err := getAuthTransactor(client, mnemonic)
	if err != nil {
		panic(err)
	}

	_, owner := getPrivateKey(mnemonic)
	ownerBalance, err := contract.BalanceOf(nil, owner)
	if err != nil {
		log.Error("cannot get balance for owner: ", owner, " err = ", err)
	}

	// Check the allowance to see if we can quit early.
	allowance, err := contract.Allowance(nil, owner, common.HexToAddress(target))
	if err != nil {
		log.Error("cannot get allowance, err = ", err)
	}
	if allowance.Cmp(ownerBalance) >= 0 {
		log.Verbose("The target has enough balance. No need to approve.")
		return
	}

	// Make a tx to approve.
	log.Verbose("Approving address ", target, " token = ", erc20Addr,
		" owner balance = ", ownerBalance, " nonce = ", opts.Nonce)

	tx, err := contract.Approve(opts, common.HexToAddress(target), ownerBalance)
	if err != nil {
		log.Error("Cannot approve address, err = ", err)
	}
	bind.WaitDeployed(context.Background(), client, tx)
	time.Sleep(time.Second * 5)
}

func transferErc20(client *ethclient.Client, mnemonic string, erc20Addr string, target string) {
	contract, err := erc20.NewErc20(common.HexToAddress(erc20Addr), client)
	if err != nil {
		panic(err)
	}

	opts, err := getAuthTransactor(client, mnemonic)
	if err != nil {
		panic(err)
	}

	tx, err := contract.Transfer(opts, common.HexToAddress(target), new(big.Int).Mul(utils.EthToWei, big.NewInt(10000)))
	if err != nil {
		log.Error("Cannot approve address, err = ", err)
	}
	bind.WaitDeployed(context.Background(), client, tx)
	time.Sleep(time.Second * 5)
}

func queryToken(ctx context.Context, sisuRpc, tokenId string) *tssTypes.Token {
	grpcConn, err := grpc.Dial(
		sisuRpc,
		grpc.WithInsecure(),
	)
	defer grpcConn.Close()
	if err != nil {
		panic(err)
	}

	queryClient := tssTypes.NewTssQueryClient(grpcConn)
	res, err := queryClient.QueryToken(context.Background(), &tssTypes.QueryTokenRequest{
		Id: tokenId,
	})
	if err != nil {
		panic(err)
	}

	return res.Token
}

func getEthVaultAddress(context context.Context, chain string, sisuRpc string) string {
	grpcConn, err := grpc.Dial(
		sisuRpc,
		grpc.WithInsecure(),
	)
	defer grpcConn.Close()
	if err != nil {
		panic(err)
	}

	queryClient := tssTypes.NewTssQueryClient(grpcConn)
	res, err := queryClient.QueryVault(context, &tssTypes.QueryVaultRequest{
		Chain: chain,
	})

	if err != nil {
		panic(err)
	}

	if len(res.Vault.Address) == 0 {
		panic("gateway contract address is empty")
	}

	return res.Vault.Address
}

type Engine struct {
	Client *ethclient.Client
	owner  common.Address
	opts   *bind.TransactOpts
}

func NewEngine(client *ethclient.Client, mnemonic string) (*Engine, error) {
	privateKey, owner := getPrivateKey(mnemonic)

	// This is the private key of the accounts0
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return nil, err
	}

	auth.Value = big.NewInt(0)
	auth.GasLimit = 6_000_000

	return &Engine{
		Client: client,
		owner:  owner,
		opts:   auth,
	}, nil
}

func (g *Engine) SetValue(value *big.Int) {
	g.opts.Value = value
}

func (g *Engine) SetGasLimit(gasLimit uint64) {
	g.opts.GasLimit = gasLimit
}

func (g *Engine) Run(f func(opts *bind.TransactOpts) *etypes.Transaction) {
	nonce, err := g.Client.PendingNonceAt(context.Background(), g.owner)
	if err != nil {
		panic(err)
	}
	gasPrice, err := g.Client.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}

	g.opts.Nonce = big.NewInt(int64(nonce))
	g.opts.GasPrice = gasPrice

	tx := f(g.opts)

	receipt, err := bind.WaitMined(context.Background(), g.Client, tx)
	if err != nil {
		panic(err)
	}

	log.Info("txHash = ", receipt.TxHash.Hex())
}

func (g *Engine) Deploy(f func(opts *bind.TransactOpts) *etypes.Transaction) common.Address {
	nonce, err := g.Client.PendingNonceAt(context.Background(), g.owner)
	if err != nil {
		panic(err)
	}
	gasPrice, err := g.Client.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}

	g.opts.Nonce = big.NewInt(int64(nonce))
	g.opts.GasPrice = gasPrice

	tx := f(g.opts)

	log.Info("Deploying contract ...")
	contractAddr, err := bind.WaitDeployed(context.Background(), g.Client, tx)
	if err != nil {
		panic(err)
	}

	log.Info("Deployed contract successfully, addr: ", contractAddr.String())
	return contractAddr
}

func getAuthTransactor(client *ethclient.Client, mnemonic string) (*bind.TransactOpts, error) {
	// This is the private key of the accounts0
	privateKey, owner := getPrivateKey(mnemonic)
	nonce, err := client.PendingNonceAt(context.Background(), owner)
	if err != nil {
		return nil, err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	// This is the private key of the accounts0
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasPrice = gasPrice

	auth.GasLimit = uint64(6_000_000)

	return auth, nil
}
