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
	UpdateContractsStatusFunc     func(contracts []*tsstypes.ContractEntity, status string)
	UpdateContractDeployTxFunc    func(chain, id string, txHash string)
	UpdateContractAddressFunc     func(chain, hash, address string)

	// Txout
	InsertTxOutsFunc       func(txs []*tsstypes.TxOutEntity)
	GetTxOutWithHashFunc   func(chain string, hash string, isHashWithSig bool) *tsstypes.TxOutEntity
	IsContractDeployTxFunc func(chain string, hashWithoutSig string) bool
	UpdateTxOutSigFunc     func(chain, hashWithoutSign, hashWithSig string, sig []byte)
	UpdateTxOutStatusFunc  func(chain, hashWithoutSig, status string)

	// Mempool tx
	InsertMempoolTxHashFunc   func(hash string, blockHeight int64)
	MempoolTxExistedFunc      func(hash string) bool
	MempoolTxExistedRangeFunc func(hash string, minBlock int64, maxBlock int64) bool
}

func (d *Database) Init() error {
	panic("implement me")
}

func (d *Database) Close() error {
	panic("implement me")
}

func (d *Database) CreateKeygen(chain string) error {
	panic("implement me")
}

func (d *Database) UpdateKeygenAddress(chain, address string, pubKey []byte) {
	panic("implement me")
}

func (d *Database) IsKeyExisted(chain string) bool {
	panic("implement me")
}

func (d *Database) IsChainKeyAddress(chain, address string) bool {
	panic("implement me")
}

func (d *Database) GetPubKey(chain string) []byte {
	panic("implement me")
}

func (d *Database) UpdateKeygenStatus(chain, status string) {
	panic("implement me")
}

func (d *Database) GetKeygenStatus(chain string) (string, error) {
	panic("implement me")
}

func (d *Database) InsertContracts(contracts []*tsstypes.ContractEntity) {
	panic("implement me")
}

func (d *Database) GetPendingDeployContracts(chain string) []*tsstypes.ContractEntity {
	panic("implement me")
}

func (d *Database) GetContractFromAddress(chain, address string) *tsstypes.ContractEntity {
	panic("implement me")
}

func (d *Database) GetContractFromHash(chain, hash string) *tsstypes.ContractEntity {
	panic("implement me")
}

func (d *Database) UpdateContractsStatus(contracts []*tsstypes.ContractEntity, status string) {
	panic("implement me")
}

func (d *Database) UpdateContractDeployTx(chain, id string, txHash string) {
	panic("implement me")
}

func (d *Database) UpdateContractAddress(chain, hash, address string) {
	panic("implement me")
}

func (d *Database) InsertTxOuts(txs []*tsstypes.TxOutEntity) {
	panic("implement me")
}

func (d *Database) GetTxOutWithHash(chain string, hash string, isHashWithSig bool) *tsstypes.TxOutEntity {
	panic("implement me")
}

func (d *Database) IsContractDeployTx(chain string, hashWithoutSig string) bool {
	panic("implement me")
}

func (d *Database) UpdateTxOutSig(chain, hashWithoutSign, hashWithSig string, sig []byte) {
	panic("implement me")
}

func (d *Database) UpdateTxOutStatus(chain, hashWithoutSig, status string) {
	panic("implement me")
}

func (d *Database) InsertMempoolTxHash(hash string, blockHeight int64) {
	panic("implement me")
}

func (d *Database) MempoolTxExisted(hash string) bool {
	panic("implement me")
}

func (d *Database) MempoolTxExistedRange(hash string, minBlock int64, maxBlock int64) bool {
	panic("implement me")
}

