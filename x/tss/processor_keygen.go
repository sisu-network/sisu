package tss

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	dhTypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/sisu/contracts/eth/dummy"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
)

/**
Process for generating a new key:
- Wait for the app to catch up
- If there is no support for a particular chain, creates a proposal to include a chain
- When other nodes receive the proposal, top N validator nodes vote to see if it should accept that.
- After M blocks (M is a constant) since a proposal is sent, count the number of yes vote. If there
are enough validator supporting the new chain, send a message to TSS engine to do keygen.
*/

type BlockSymbolPair struct {
	blockHeight int64
	chainSymbol string
}

func (p *Processor) CheckTssKeygen(ctx sdk.Context, blockHeight int64) {
	if p.globalData.IsCatchingUp() ||
		p.lastProposeBlockHeight != 0 && blockHeight-p.lastProposeBlockHeight <= PROPOSE_BLOCK_INTERVAL {
		return
	}

	chains, err := p.keeper.GetRecordedChainsOnSisu(ctx)
	if err != nil {
		return
	}

	utils.LogInfo("recordedChains = ", chains)

	unavailableChains := make([]string, 0)
	for _, chainConfig := range p.config.SupportedChains {
		// TODO: Remove this after testing.
		// if chains.Chains[chainConfig.Symbol] == nil {
		if true {
			unavailableChains = append(unavailableChains, chainConfig.Symbol)
		}
	}

	// Broadcast a message.
	utils.LogInfo("Broadcasting TSS Keygen Proposal message. len(unavailableChains) = ", len(unavailableChains))
	signer := p.appKeys.GetSignerAddress()

	for _, chain := range unavailableChains {
		proposal := types.NewMsgKeygenProposal(
			signer.String(),
			chain,
			utils.GenerateRandomString(16),
			blockHeight+int64(p.config.BlockProposalLength),
		)
		utils.LogDebug("Submitting proposal message for chain", chain)
		go func() {
			err := p.txSubmit.SubmitMessage(proposal)
			if err != nil {
				utils.LogError(err)
			}
		}()
	}

	p.lastProposeBlockHeight = blockHeight
}

// Called after having key generation result from Sisu's api server.
func (p *Processor) OnKeygenResult(result dhTypes.KeygenResult) {
	// 1. Post result to the cosmos chain
	signer := p.appKeys.GetSignerAddress()

	resultEnum := types.KeygenResult_FAILURE
	if result.Success {
		resultEnum = types.KeygenResult_SUCCESS
	}

	msg := types.NewKeygenResult(signer.String(), result.Chain, resultEnum, result.PubKeyBytes, result.Address)
	p.txSubmit.SubmitMessage(msg)

	// 2. Add the address to the watch list.
	deyesClient := p.deyesClients[result.Chain]
	if deyesClient == nil {
		utils.LogCritical("Cannot find deyes client for chain", result.Chain)
	} else {
		deyesClient.AddWatchAddresses(result.Chain, []string{result.Address})
	}
}

func (p *Processor) CheckKeyGenProposal(msg *types.KeygenProposal) error {
	// TODO: Check if we see the same need to have keygen proposal here.
	return nil
}

func (p *Processor) DeliverKeyGenProposal(msg *types.KeygenProposal) ([]byte, error) {
	// Send a signal to Dheart to start keygen process.
	utils.LogInfo("Sending keygen request to Dheart...")
	pubKeys := p.partyManager.GetActivePartyPubkeys()
	keygenId := GetKeygenId(msg.ChainSymbol, p.currentHeight, pubKeys)
	err := p.dheartClient.KeyGen(keygenId, msg.ChainSymbol, pubKeys)
	if err != nil {
		utils.LogError(err)
		return nil, err
	}
	utils.LogInfo("Keygen request is sent successfully.")

	return []byte{}, nil
}

func (p *Processor) DeliverKeygenResult(ctx sdk.Context, msg *types.KeygenResult) ([]byte, error) {
	// TODO: Accumulates results from others and check for bad actors

	// For now, only process self message sent from this node.
	if msg.Signer == p.appKeys.GetSignerAddress().String() {
		utils.LogDebug("Keygen: This is the same signer...")

		// Save this to KVStore
		chainsInfo, err := p.keeper.GetRecordedChainsOnSisu(ctx)
		if err != nil {
			utils.LogError(err)
			return nil, err
		}

		if chainsInfo.Chains == nil {
			chainsInfo.Chains = make(map[string]*types.ChainInfo)
		}

		// TODO: Add validators here.
		chainsInfo.Chains[msg.ChainSymbol] = &types.ChainInfo{
			Symbol: msg.ChainSymbol,
		}

		p.keeper.SetChainsInfo(ctx, chainsInfo)

		// Save the pubkey to the keeper.
		p.keeper.SavePubKey(ctx, msg.ChainSymbol, msg.PubKeyBytes)

		// If this is a pubkey address of a ETH chain, save it to the store because we want to watch
		// transaction that funds the address (we will deploy contracts later).
		if utils.IsETHBasedChain(msg.ChainSymbol) {
			p.txOutputProducer.AddKeyAddress(ctx, msg.ChainSymbol, msg.Address)
		}

		// Check and see if we need to deploy some contracts. If we do, push them into the contract
		// queue for deployment later (after we receive some funding like ether to execute contract
		// deployment).
		p.checkContractDeployment(ctx, msg)

		p.printKeygenPubKey(msg)
	} else {
		utils.LogDebug("Keygen: message is from different signers.")
	}

	return nil, nil
}

// Print out the public key address. Used for debugging purpose
func (p *Processor) printKeygenPubKey(msg *types.KeygenResult) {
	pubKey, err := crypto.DecompressPubkey(msg.PubKeyBytes)
	if err == nil {
		// TODO: Check if the chain is ETH before getting public key.
		address := crypto.PubkeyToAddress(*pubKey).Hex()
		utils.LogInfo("Address = ", address)
	} else {
		utils.LogError("Critical Error, public key cannot be deserialized. Err = ", err)
	}
}

// Called after keygen finishes to see if we need to deploy any contracts.
func (p *Processor) checkContractDeployment(ctx sdk.Context, msg *types.KeygenResult) {
	contractABIs := []string{
		dummy.DummyABI,
	}

	for _, abi := range contractABIs {
		// Hash of a contract is the hash of the ABI string.
		hash, err := utils.KeccakHash32(abi)
		if err != nil {
			utils.LogError("Cannot get keccak hash 32 byte, err = ", err)
			continue
		}

		// Check if this contract has been deployed or being deployed.
		if p.keeper.IsContractDeployingOrDeployed(ctx, msg.ChainSymbol, hash) {
			utils.LogDebug("Contract has been deployed or being deployed. Hash = ", hash)
			continue
		}

		p.keeper.EnqueueContract(ctx, msg.ChainSymbol, hash, abi)
	}
}
