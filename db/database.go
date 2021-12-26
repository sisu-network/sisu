package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/tss/types"
	tsstypes "github.com/sisu-network/sisu/x/tss/types"
)

//go:generate mockgen -source=db/database.go -destination=tests/mock/db.go -package=mock

// Make sure struct implement interface at compile-time
var _ Database = (*SqlDatabase)(nil)

// This is an interface of a private database that is only be used by this node and not other nodes.
type Database interface {
	Init() error
	Close() error

	// Keygen
	CreateKeygen(keyType string, startBlock int64) error
	GetKeyGen(keyType string) (*tsstypes.KeygenEntity, error)

	UpdateKeygenAddress(keyType, address string, pubKey []byte)
	IsChainKeyAddress(keyType string, address string) bool
	GetPubKey(keyType string) []byte
	UpdateKeygenStatus(keyType, status string)

	// Contracts
	GetContractFromAddress(chain, address string) *tsstypes.ContractEntity
	GetContractFromHash(chain, hash string) *tsstypes.ContractEntity
	UpdateContractsStatus(contracts []*tsstypes.ContractEntity, status string) error
	UpdateContractDeployTx(chain, id string, txHash string)
	UpdateContractAddress(chain, hash, address string)

	// TxIn
	InsertTxIn(txIn *types.TxIn)
	IsTxInExisted(txIn *types.TxIn) bool

	// Txout
	InsertTxOuts(txs []*types.TxOut)
	GetTxOutWithHash(chain string, hash string, isHashWithSig bool) *tsstypes.TxOutEntity
	IsContractDeployTx(chain string, hashWithoutSig string) bool
	UpdateTxOutSig(chain, hashWithoutSign, hashWithSig string, sig []byte) error
	UpdateTxOutStatus(chain, hash string, status tsstypes.TxOutStatus, isHashWithSig bool) error

	// Mempool tx
	InsertMempoolTxHash(hash string, blockHeight int64)
	MempoolTxExisted(hash string) bool
	MempoolTxExistedRange(hash string, minBlock int64, maxBlock int64) bool
}

type SqlDatabase struct {
	db     *sql.DB
	config config.SqlConfig
}

type dbLogger struct {
}

func (loggger *dbLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (loggger *dbLogger) Verbose() bool {
	return true
}

func NewDatabase(config config.SqlConfig) Database {
	return &SqlDatabase{
		config: config,
	}
}

func (d *SqlDatabase) Init() error {
	err := d.Connect()
	if err != nil {
		log.Error("Failed to connect to DB. Err =", err)
		return err
	}

	err = d.DoMigration()
	if err != nil {
		log.Error("Cannot do migration. Err =", err)
		return err
	}

	return nil
}

func (d *SqlDatabase) Connect() error {
	host := d.config.Host
	if host == "" {
		return fmt.Errorf("DB host cannot be empty")
	}

	username := d.config.Username
	password := d.config.Password
	schema := d.config.Schema

	// Connect to the db
	database, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/", username, password, host, d.config.Port))
	if err != nil {
		return err
	}
	_, err = database.Exec("CREATE DATABASE IF NOT EXISTS " + schema)
	if err != nil {
		return err
	}
	database.Close()

	database, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, d.config.Port, schema))
	if err != nil {
		return err
	}

	d.db = database
	log.Info("Db is connected successfully")
	return nil
}

func (d *SqlDatabase) DoMigration() error {
	driver, err := mysql.WithInstance(d.db, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations/",
		"mysql",
		driver,
	)

	if err != nil {
		return err
	}

	log.Info("Doing sql migration...")

	m.Log = &dbLogger{}
	m.Up()

	return nil
}

func (d *SqlDatabase) Close() error {
	return d.db.Close()
}

func (d *SqlDatabase) CreateKeygen(keyType string, startBlock int64) error {
	query := "INSERT INTO keygen (key_type, start_block, status) VALUES (?, ?, ?)"
	params := []interface{}{keyType, startBlock, types.KEYGEN_STATUS_GENERATING}

	_, err := d.db.Exec(query, params...)
	if err != nil {
		log.Error("failed to create new keygen with type", keyType, ", err = ", err)
		return err
	}

	return nil
}

func (d *SqlDatabase) GetKeyGen(keyType string) (*tsstypes.KeygenEntity, error) {
	query := "SELECT key_type, address, pubkey, status, start_block FROM keygen WHERE key_type = ?"
	params := []interface{}{keyType}

	rows, err := d.db.Query(query, params...)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		var nullableType, nullableAddress, nullableStatus sql.NullString
		result := new(tsstypes.KeygenEntity)
		if err := rows.Scan(&nullableType, &nullableAddress, &result.Pubkey, &nullableStatus, &result.StartBlock); err != nil {
			return nil, err
		}
		result.Type = nullableType.String
		result.Address = nullableAddress.String
		result.Status = nullableStatus.String

		return result, nil
	}

	return nil, nil
}

func (d *SqlDatabase) UpdateKeygenAddress(keyType, address string, pubKey []byte) {
	query := "UPDATE keygen SET address = ?, pubkey = ? WHERE key_type = ?"
	params := []interface{}{address, pubKey, keyType}

	_, err := d.db.Exec(query, params...)
	if err != nil {
		log.Error("failed to update keygen address and pubkey, err = ", err)
	}
}

func (d *SqlDatabase) IsKeyExisted(keyType string) bool {
	query := "SELECT key_type FROM keygen WHERE key_type = ?"
	params := []interface{}{keyType}

	rows, err := d.db.Query(query, params...)
	if err != nil {
		log.Error("cannot query chain key ", keyType)
		return false
	}
	defer rows.Close()

	return rows.Next()
}

func (d *SqlDatabase) IsChainKeyAddress(keyType, address string) bool {
	query := "SELECT address FROM keygen WHERE key_type = ? AND address = ?"
	params := []interface{}{keyType, address}

	rows, err := d.db.Query(query, params...)
	if err != nil {
		log.Error("cannot query chain key ", address)
		return false
	}

	defer rows.Close()

	return rows.Next()
}

func (d *SqlDatabase) GetPubKey(keyType string) []byte {
	query := "SELECT pubkey FROM keygen WHERE key_type = ?"
	params := []interface{}{keyType}

	rows, err := d.db.Query(query, params...)
	if err != nil {
		log.Error("cannot query pub key", keyType)
		return nil
	}
	defer rows.Close()

	if !rows.Next() {
		return nil
	}

	var result []byte
	if err := rows.Scan(&result); err != nil {
		log.Error("cannot scan result, err = ", err)
		return nil
	}

	return result
}

func (d *SqlDatabase) UpdateKeygenStatus(keyType, status string) {
	query := "UPDATE keygen SET status = ? WHERE key_type = ?"
	params := []interface{}{status, keyType}

	_, err := d.db.Exec(query, params...)
	if err != nil {
		log.Error("failed to udpate keygen status for key type", keyType, ", err = ", err)
	}
}
func (d *SqlDatabase) GetContractFromAddress(chain, address string) *tsstypes.ContractEntity {
	query := "SELECT chain, hash, byteCode, name, address, status FROM contract WHERE chain=? AND address = ?"
	params := []interface{}{chain, address}

	rows, err := d.db.Query(query, params...)
	if err != nil {
		return nil
	}

	defer rows.Close()

	if rows.Next() {
		var chain, hash, name, address, status sql.NullString
		var byteCode []byte

		if err := rows.Scan(&chain, &hash, &byteCode, &name, &address, &status); err != nil {
			return nil
		}

		return &tsstypes.ContractEntity{
			Chain:   chain.String,
			Hash:    hash.String,
			Name:    name.String,
			Address: address.String,
			Status:  status.String,
		}
	}

	return nil
}

func (d *SqlDatabase) GetContractFromHash(chain, hash string) *tsstypes.ContractEntity {
	query := "SELECT chain, hash, byteCode, name, address, status FROM contract WHERE chain=? AND hash = ?"
	params := []interface{}{chain, hash}

	rows, err := d.db.Query(query, params...)
	if err != nil {
		return nil
	}

	defer rows.Close()

	if rows.Next() {
		var chain, hash, name, address, status sql.NullString
		var byteCode []byte

		if err := rows.Scan(&chain, &hash, &byteCode, &name, &address, &status); err != nil {
			return nil
		}

		return &tsstypes.ContractEntity{
			Chain:   chain.String,
			Hash:    hash.String,
			Name:    name.String,
			Address: address.String,
			Status:  status.String,
		}
	}

	return nil
}

func (d *SqlDatabase) UpdateContractsStatus(contracts []*tsstypes.ContractEntity, status string) error {
	for _, contract := range contracts {
		query := "UPDATE contract SET status = ? WHERE chain = ? AND hash = ?"
		params := []interface{}{status, contract.Chain, contract.Hash}

		if _, err := d.db.Exec(query, params...); err != nil {
			log.Error("failed to update contract status, err =", err, ". len(contracts) =", len(contracts))
			return err
		}
	}

	return nil
}

func (d *SqlDatabase) UpdateContractDeployTx(chain, hash string, txHash string) {
	query := "UPDATE contract SET tx_hash = ? WHERE chain = ? AND hash = ?"
	params := []interface{}{txHash, chain, hash}

	_, err := d.db.Exec(query, params...)
	if err != nil {
		log.Error("failed to update contract deploy tx, err =", err)
	}
}

func (d *SqlDatabase) IsContractDeployTx(chain string, hashWithoutSig string) bool {
	query := "SELECT contract_hash FROM tx_out WHERE chain=? AND hash_without_sig = ?"
	params := []interface{}{
		chain,
		hashWithoutSig,
	}

	rows, err := d.db.Query(query, params...)
	if err != nil {
		return false
	}

	defer rows.Close()

	if rows.Next() {
		var hash sql.NullString
		if err := rows.Scan(&hash); err != nil {
			return false
		}

		return hash.String != ""
	}

	return false
}

func (d *SqlDatabase) UpdateContractAddress(chain, outHash, address string) {
	query := "UPDATE contract SET address = ? WHERE chain = ? AND hash = (SELECT contract_hash FROM tx_out WHERE chain = ? AND hash_without_sig = ?)"
	params := []interface{}{address, chain, chain, outHash}

	_, err := d.db.Exec(query, params...)
	if err != nil {
		log.Error("failed to update contract address, err =", err)
	}
}

func (d *SqlDatabase) InsertTxIn(txIn *types.TxIn) {
	query := "INSERT IGNORE INTO tx_in (chain, hash, block_height, serialized) VALUES(?, ?, ?, ?)"
	params := []interface{}{txIn.Chain, txIn.TxHash, txIn.BlockHeight, txIn.Serialized}

	_, err := d.db.Exec(query, params...)
	if err != nil {
		log.Error("failed to insert TxIn into table, err =", err)
	}
}

func (d *SqlDatabase) IsTxInExisted(txIn *types.TxIn) bool {
	query := "SELECT chain FROM tx_in WHERE chain=? AND hash=? AND block_height=?"
	params := []interface{}{txIn.Chain, txIn.TxHash, txIn.BlockHeight}

	rows, err := d.db.Query(query, params...)
	if err != nil {
		log.Error("failed to query TxIn, err =", err)
		return false
	}

	defer rows.Close()

	return rows.Next()
}

func (d *SqlDatabase) InsertTxOuts(txs []*types.TxOut) {
	query := "INSERT INTO tx_out (chain, hash_without_sig, in_chain, in_hash, bytes_without_sig) VALUES "
	query = query + getQueryQuestionMark(len(txs), 5)

	params := make([]interface{}, 0, len(txs)*6)

	for _, tx := range txs {
		params = append(params, tx.OutChain)
		params = append(params, tx.GetHash())
		params = append(params, tx.InChain)
		params = append(params, tx.InHash)
		params = append(params, tx.OutBytes)
	}

	_, err := d.db.Exec(query, params...)
	if err != nil {
		log.Error("failed to insert txout into table, err = ", err)
	}
}

func (d *SqlDatabase) GetTxOutWithHash(chain string, hash string, isHashWithSig bool) *tsstypes.TxOutEntity {
	var query string
	if isHashWithSig {
		query = "SELECT chain, status, hash_without_sig, hash_with_sig, in_chain, in_hash, bytes_without_sig, signature, contract_hash FROM tx_out WHERE chain = ? AND hash_with_sig = ?"
	} else {
		query = "SELECT chain, status, hash_without_sig, hash_with_sig, in_chain, in_hash, bytes_without_sig, signature, contract_hash FROM tx_out WHERE chain = ? AND hash_without_sig = ?"
	}
	params := []interface{}{chain, hash}

	rows, err := d.db.Query(query, params...)
	if err != nil {
		return nil
	}

	defer rows.Close()

	if rows.Next() {
		var chain, status, hashWithoutSig, hashWithSig, inChain, inHash, contractHash sql.NullString
		var bytesWithoutSig, signature []byte

		if err := rows.Scan(&chain, &status, &hashWithoutSig, &hashWithSig, &inChain, &inHash, &bytesWithoutSig, &signature, &contractHash); err != nil {
			return nil
		}

		return &tsstypes.TxOutEntity{
			OutChain:        chain.String,
			HashWithoutSig:  hashWithoutSig.String,
			HashWithSig:     hashWithSig.String,
			InChain:         inChain.String,
			InHash:          inHash.String,
			BytesWithoutSig: bytesWithoutSig,
			Status:          status.String,
			Signature:       string(signature),
		}
	}

	return nil
}

func (d *SqlDatabase) UpdateTxOutSig(chain, hashWithoutSign, hashWithSig string, sig []byte) error {
	query := "UPDATE tx_out SET signature = ?, hash_with_sig = ? WHERE chain = ? AND hash_without_sig = ?"
	params := []interface{}{
		sig,
		hashWithSig,
		chain,
		hashWithoutSign,
	}

	if _, err := d.db.Exec(query, params...); err != nil {
		log.Error("failed to update txout with chain and hashWoSig", chain, hashWithSig, ", err =", err)
		return err
	}

	return nil
}

func (d *SqlDatabase) UpdateTxOutStatus(chain, hash string, status tsstypes.TxOutStatus, isHashWithSig bool) error {
	log.Debugf("Updating txout hash(%s) to status(%s), chain(%s)", hash, string(status), chain)
	query := "UPDATE tx_out SET status = ? WHERE chain = ? AND hash_with_sig = ?"
	if !isHashWithSig {
		query = "UPDATE tx_out SET status = ? WHERE chain = ? AND hash_without_sig = ?"
	}

	params := []interface{}{status, chain, hash}

	if _, err := d.db.Exec(query, params...); err != nil {
		log.Error("failed to update chain status", chain, hash, ", err =", err)
		return err
	}

	return nil
}

func (d *SqlDatabase) InsertMempoolTxHash(hash string, blockHeight int64) {
	query := "INSERT INTO mempool_tx (hash, block_height) VALUES (?, ?)"
	params := []interface{}{hash, blockHeight}

	_, err := d.db.Exec(query, params...)
	if err != nil {
		log.Error("failed to insert tx hash into mempool_tx table, err =", err)
	}
}

func (d *SqlDatabase) MempoolTxExisted(hash string) bool {
	query := "SELECT hash FROM mempool_tx WHERE hash=?"
	params := []interface{}{hash}

	rows, err := d.db.Query(query, params...)
	if err != nil {
		log.Error("failed to query mempool_tx, err =", err)
	}
	defer rows.Close()

	return rows.Next()
}

func (d *SqlDatabase) MempoolTxExistedRange(hash string, minBlock int64, maxBlock int64) bool {
	query := "SELECT hash FROM mempool_tx WHERE hash=? AND block_height >= ? AND block_height <= ?"
	params := []interface{}{hash, minBlock, maxBlock}

	rows, err := d.db.Query(query, params...)
	if err != nil {
		log.Error("failed to query mempool_tx, err =", err)
	}
	defer rows.Close()

	return rows.Next()
}
