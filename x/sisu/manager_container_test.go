package sisu

import (
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
)

func MockManagerContainer(args ...interface{}) ManagerContainer {
	mc := &DefaultManagerContainer{}

	for _, arg := range args {
		switch t := arg.(type) {
		case PostedMessageManager:
			mc.pmm = t
		case keeper.Storage:
			mc.publicDb = t
		case common.GlobalData:
			mc.globalData = t
		case tssclients.DeyesClient:
			mc.deyesClient = t
		case tssclients.DheartClient:
			mc.dheartClient = t
		case config.TssConfig:
			mc.config = t
		case common.TxSubmit:
			mc.txSubmit = t
		case common.AppKeys:
			mc.appKeys = t
		case PartyManager:
			mc.partyManager = t
		}
	}

	return mc
}
