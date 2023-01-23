package background

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/chains"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/keeper"
)

///// ManagerContainer

func MockManagerContainer(args ...interface{}) ManagerContainer {
	mc := &DefaultManagerContainer{}

	for _, arg := range args {
		switch t := arg.(type) {
		case components.PostedMessageManager:
			mc.pmm = t
		case components.GlobalData:
			mc.globalData = t
		case external.DeyesClient:
			mc.deyesClient = t
		case external.DheartClient:
			mc.dheartClient = t
		case config.Config:
			mc.config = t
		case components.TxSubmit:
			mc.txSubmit = t
		case components.AppKeys:
			mc.appKeys = t
		case components.PartyManager:
			mc.partyManager = t
		case components.TxTracker:
			mc.txTracker = t
		case keeper.Keeper:
			mc.keeper = t
		case sdk.Context:
			mc.readOnlyContext.Store(t)
		case chains.TxOutputProducer:
			mc.txOutProducer = t
		case components.ValidatorManager:
			mc.valsManager = t
		case TransferQueue:
			mc.transferOutQueue = t
		case chains.BridgeManager:
			mc.bridgeManager = t
		case keeper.PrivateDb:
			mc.privateDb = t
		}
	}

	return mc
}

