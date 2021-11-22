package mock

import (
	"github.com/sisu-network/sisu/db"
	tsstypes "github.com/sisu-network/sisu/x/tss/types"
)

var _ db.Database = (*Database)(nil)

type Database struct {
	InitFunc  func() error
	CloseFunc func() error

	// Keygen
	CreateKeygenFunc        func(chain string) error
	UpdateKeygenAddressFunc func(chain, address string, pubKey []byte)

	IsKeyExistedFunc       func(chain string) bool
	IsChainKeyAddressFunc  func(chain, address string) bool
	GetPubKeyFunc          func(chain string) []byte
	UpdateKeygenStatusFunc func(chain, status string)
	GetKeygenStatusFunc    func(chain string) (string, error)

	// Contracts
	InsertContractsFunc           func(contracts []*tsstypes.ContractEntity)
	GetPendingDeployContractsFunc func(chain string) []*tsstypes.ContractEntity
	GetContractFromAddressFunc    func(chain, address string) *tsstypes.ContractEntity
	GetContractFromHashFunc       func(chain, hash string) *tsstypes.ContractEntity
	UpdateContractsStatusFunc     func(contracts []*tsstypes.ContractEntity, status string) error
	UpdateContractDeployTxFunc    func(chain, id string, txHash string)
	UpdateContractAddressFunc     func(chain, hash, address string)

	// Txout
	InsertTxOutsFunc       func(txs []*tsstypes.TxOutEntity)
	GetTxOutWithHashFunc   func(chain string, hash string, isHashWithSig bool) *tsstypes.TxOutEntity
	IsContractDeployTxFunc func(chain string, hashWithoutSig string) bool
	UpdateTxOutSigFunc     func(chain, hashWithoutSign, hashWithSig string, sig []byte) error
	UpdateTxOutStatusFunc  func(chain, hashWithoutSig, status string) error

	// Mempool tx
	InsertMempoolTxHashFunc   func(hash string, blockHeight int64)
	MempoolTxExistedFunc      func(hash string) bool
	MempoolTxExistedRangeFunc func(hash string, minBlock int64, maxBlock int64) bool
}

func (d *Database) Init() error {
	if d.InitFunc == nil {
		panic("function is not defined")
	}

	return d.InitFunc()
}

func (d *Database) Close() error {
	if d.CloseFunc == nil {
		panic("function is not defined")
	}

	return d.CloseFunc()
}

func (d *Database) CreateKeygen(chain string) error {
	if d.CreateKeygenFunc == nil {
		panic("function is not defined")
	}

	return d.CreateKeygenFunc(chain)
}

func (d *Database) UpdateKeygenAddress(chain, address string, pubKey []byte) {
	if d.UpdateKeygenAddressFunc == nil {
		panic("function is not defined")
	}

	d.UpdateKeygenAddressFunc(chain, address, pubKey)
}

func (d *Database) IsKeyExisted(chain string) bool {
	if d.IsKeyExistedFunc == nil {
		panic("function is not defined")
	}

	return d.IsKeyExistedFunc(chain)
}

func (d *Database) IsChainKeyAddress(chain, address string) bool {
	if d.IsChainKeyAddressFunc == nil {
		panic("function is not defined")
	}

	return d.IsChainKeyAddressFunc(chain, address)
}

func (d *Database) GetPubKey(chain string) []byte {
	if d.GetPubKeyFunc == nil {
		panic("function is not defined")
	}

	return d.GetPubKeyFunc(chain)
}

func (d *Database) UpdateKeygenStatus(chain, status string) {
	if d.UpdateKeygenStatusFunc == nil {
		panic("function is not defined")
	}

	d.UpdateKeygenStatusFunc(chain, status)
}

func (d *Database) GetKeygenStatus(chain string) (string, error) {
	if d.GetKeygenStatusFunc == nil {
		panic("function is not defined")
	}

	return d.GetKeygenStatusFunc(chain)
}

func (d *Database) InsertContracts(contracts []*tsstypes.ContractEntity) {
	if d.InsertContractsFunc == nil {
		panic("function is not defined")
	}

	d.InsertContractsFunc(contracts)
}

func (d *Database) GetPendingDeployContracts(chain string) []*tsstypes.ContractEntity {
	if d.GetPendingDeployContractsFunc == nil {
		panic("function is not defined")
	}

	return d.GetPendingDeployContractsFunc(chain)
}

func (d *Database) GetContractFromAddress(chain, address string) *tsstypes.ContractEntity {
	if d.GetContractFromAddressFunc == nil {
		panic("function is not defined")
	}

	return d.GetContractFromAddressFunc(chain, address)
}

func (d *Database) GetContractFromHash(chain, hash string) *tsstypes.ContractEntity {
	if d.GetContractFromHashFunc == nil {
		panic("function is not defined")
	}

	return d.GetContractFromHashFunc(chain, hash)
}

func (d *Database) UpdateContractsStatus(contracts []*tsstypes.ContractEntity, status string) error {
	if d.UpdateContractsStatusFunc == nil {
		panic("function is not defined")
	}

	return d.UpdateContractsStatusFunc(contracts, status)
}

func (d *Database) UpdateContractDeployTx(chain, id string, txHash string) {
	if d.UpdateContractDeployTxFunc == nil {
		panic("function is not defined")
	}

	d.UpdateContractDeployTxFunc(chain, id, txHash)
}

func (d *Database) UpdateContractAddress(chain, hash, address string) {
	if d.UpdateContractAddressFunc == nil {
		panic("function is not defined")
	}

	d.UpdateContractAddressFunc(chain, hash, address)
}

func (d *Database) InsertTxOuts(txs []*tsstypes.TxOutEntity) {
	if d.InsertTxOutsFunc == nil {
		panic("function is not defined")
	}

	d.InsertTxOutsFunc(txs)
}

func (d *Database) GetTxOutWithHash(chain string, hash string, isHashWithSig bool) *tsstypes.TxOutEntity {
	if d.GetTxOutWithHashFunc == nil {
		panic("function is not defined")
	}

	return d.GetTxOutWithHashFunc(chain, hash, isHashWithSig)
}

func (d *Database) IsContractDeployTx(chain string, hashWithoutSig string) bool {
	if d.IsContractDeployTxFunc == nil {
		panic("function is not defined")
	}

	return d.IsContractDeployTxFunc(chain, hashWithoutSig)
}

func (d *Database) UpdateTxOutSig(chain, hashWithoutSign, hashWithSig string, sig []byte) error {
	if d.UpdateTxOutSigFunc == nil {
		panic("function is not defined")
	}

	return d.UpdateTxOutSigFunc(chain, hashWithoutSign, hashWithSig, sig)
}

func (d *Database) UpdateTxOutStatus(chain, hashWithoutSig, status string) error {
	if d.UpdateTxOutStatusFunc == nil {
		panic("function is not defined")
	}

	return d.UpdateTxOutStatusFunc(chain, hashWithoutSig, status)
}

func (d *Database) InsertMempoolTxHash(hash string, blockHeight int64) {
	if d.InsertContractsFunc == nil {
		panic("function is not defined")
	}

	d.InsertMempoolTxHashFunc(hash, blockHeight)
}

func (d *Database) MempoolTxExisted(hash string) bool {
	if d.MempoolTxExistedFunc == nil {
		panic("function is not defined")
	}

	return d.MempoolTxExistedFunc(hash)
}

func (d *Database) MempoolTxExistedRange(hash string, minBlock int64, maxBlock int64) bool {
	if d.MempoolTxExistedRangeFunc == nil {
		panic("function is not defined")
	}

	return d.MempoolTxExistedRangeFunc(hash, minBlock, maxBlock)
}
