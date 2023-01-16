package dev

import (
	"github.com/sisu-network/deyes/config"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/helper"
)

func (c *fundAccountCmd) fundLisk(genesisFolder, mnemonic string, mpcPubKey []byte) {
	liskConfig := helper.ReadLiskConfig(genesisFolder)
	configFormatted := config.Chain{Chain: liskConfig.Chain, Rpcs: []string{liskConfig.RPC}}
	client := lisk.NewLiskRPC(configFormatted)
	mpcAddr := crypto.GetAddressFromPublicKey(mpcPubKey)
	log.Verbose("Funding LSK for mpc address = ", mpcAddr)
	amount := uint64(20000000)
	moduleId := uint32(2)
	assetId := uint32(0)
	transferLisk(client, mnemonic, mpcAddr, amount, moduleId, assetId, liskConfig)
}
