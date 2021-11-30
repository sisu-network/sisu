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
	tsstypes "github.com/sisu-network/sisu/x/tss/types"
)

//go:generate mockgen -source=database.go -destination=../tests/mock/db.go -package=mock

const (
	// The keygen proposal has passed the consensus and included in a Sisu block
	StatusProposalFinalized = "proposal_finalized"
	// The keygen has finished and delivered to destination chain.
	StatusDeliveredToChain = "delivered_to_chain"
)

// Make sure struct implement interface at compile-time
var _ Database = (*SqlDatabase)(nil)

type Database interface {
	Init() error
	Close() error

	// Keygen
	CreateKeygen(chain string) error
	UpdateKeygenAddress(chain, address string, pubKey []byte)

	IsKeyExisted(chain string) bool
	IsChainKeyAddress(chain, address string) bool
	GetPubKey(chain string) []byte
	UpdateKeygenStatus(chain, status string)
	GetKeygenStatus(chain string) (string, error)

	// Contracts
	InsertContracts(contracts []*tsstypes.ContractEntity)
	GetPendingDeployContracts(chain string) []*tsstypes.ContractEntity
	GetContractFromAddress(chain, address string) *tsstypes.ContractEntity
	GetContractFromHash(chain, hash string) *tsstypes.ContractEntity
	UpdateContractsStatus(contracts []*tsstypes.ContractEntity, status string) error
	UpdateContractDeployTx(chain, id string, txHash string)
	UpdateContractAddress(chain, hash, address string)

	// Txout
	InsertTxOuts(txs []*tsstypes.TxOutEntity)
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

func (d *SqlDatabase) CreateKeygen(chain string) error {
	query := "INSERT INTO keygen (chain) VALUES (?)"
	params := []interface{}{chain}

	_, err := d.db.Exec(query, params...)
	if err != nil {
		log.Error("failed to create new keygen for chain", chain, ", err = ", err)
		return err
	}

	return nil
}

func (d *SqlDatabase) UpdateKeygenAddress(chain, address string, pubKey []byte) {
	query := "UPDATE keygen SET address = ?, pubkey = ? WHERE chain = ?"
	params := []interface{}{address, pubKey, chain}

	_, err := d.db.Exec(query, params...)
	if err != nil {
		log.Error("failed to update keygen address and pubkey, err = ", err)
	}
}

func (d *SqlDatabase) IsKeyExisted(chain string) bool {
	query := "SELECT chain FROM keygen WHERE chain = ?"
	params := []interface{}{chain}

	rows, err := d.db.Query(query, params...)
	if err != nil {
		log.Error("cannot query chain key ", chain)
		return false
	}
	defer rows.Close()

	return rows.Next()
}

func (d *SqlDatabase) IsChainKeyAddress(chain, address string) bool {
	query := "SELECT chain FROM keygen WHERE chain = ? AND address = ?"
	params := []interface{}{chain, address}

	rows, err := d.db.Query(query, params...)
	if err != nil {
		log.Error("cannot query chain key ", chain, address)
		return false
	}

	defer rows.Close()

	return rows.Next()
}

func (d *SqlDatabase) GetPubKey(chain string) []byte {
	query := "SELECT pubkey FROM keygen WHERE chain = ?"
	params := []interface{}{chain}

	rows, err := d.db.Query(query, params...)
	if err != nil {
		log.Error("cannot query pub key", chain)
		return nil
	}
	defer rows.Close()

	if !rows.Next() {
		return nil
	}

	var result []byte
	if err := rows.Scan(&result); err != nil {
		return nil
	}

	return result
}

func (d *SqlDatabase) UpdateKeygenStatus(chain, status string) {
	query := "UPDATE keygen SET status = ? WHERE chain = ?"
	params := []interface{}{status, chain}

	_, err := d.db.Exec(query, params...)
	if err != nil {
		log.Error("failed to udpate keygen status for chain", chain, ", err = ", err)
	}
}

func (d *SqlDatabase) GetKeygenStatus(chain string) (string, error) {
	query := "SELECT status FROM keygen WHERE chain = ?"
	params := []interface{}{chain}

	rows, err := d.db.Query(query, params...)
	if err != nil {
		log.Error("cannot query keygen status ", chain)
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "", nil
	}

	var status string
	if err := rows.Scan(&status); err != nil {
		return "", err
	}

	return status, nil
}

func (d *SqlDatabase) IsKeygenDelivered(chain string) bool {
	query := "SELECT status FROM keygen WHERE chain = ?"
	params := []interface{}{chain}

	rows, err := d.db.Query(query, params...)
	if err != nil {
		log.Error("cannot check if keygen is delivered for chain ", chain, ", err =", err)
		return false
	}
	defer rows.Close()

	if !rows.Next() {
		return false
	}

	var status string
	if err := rows.Scan(&status); err != nil {
		return false
	}

	return status == StatusDeliveredToChain
}

func (d *SqlDatabase) InsertContracts(contracts []*tsstypes.ContractEntity) {
	query := "INSERT INTO contract (chain, hash, byteCode, name) VALUES "
	query = query + getQueryQuestionMark(len(contracts), 4)

	params := make([]interface{}, 0, 4*len(contracts))
	for _, contract := range contracts {
		params = append(params, contract.Chain)
		params = append(params, contract.Hash)
		params = append(params, contract.ByteCode)
		params = append(params, contract.Name)
	}

	_, err := d.db.Exec(query, params...)
	if err != nil {
		log.Error("failed to insert contract into db, err = ", err)
	}
}

func (d *SqlDatabase) GetPendingDeployContracts(chain string) []*tsstypes.ContractEntity {
	query := "SELECT chain, hash, name, status FROM contract WHERE chain=?"
	params := []interface{}{chain}
	result := make([]*tsstypes.ContractEntity, 0)

	rows, err := d.db.Query(query, params...)
	if err != nil {
		return result
	}

	defer rows.Close()

	for rows.Next() {
		var chain, hash, name, status sql.NullString
		if err := rows.Scan(&chain, &hash, &name, &status); err != nil {
			log.Error("cannot scan row, err =", err)
			continue
		}

		if status.String == "" {
			result = append(result, &tsstypes.ContractEntity{
				Chain:  chain.String,
				Hash:   hash.String,
				Name:   name.String,
				Status: status.String,
			})
		}
	}

	return result
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
			Chain:    chain.String,
			Hash:     hash.String,
			ByteCode: byteCode,
			Name:     name.String,
			Address:  address.String,
			Status:   status.String,
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
			Chain:    chain.String,
			Hash:     hash.String,
			ByteCode: byteCode,
			Name:     name.String,
			Address:  address.String,
			Status:   status.String,
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

func (d *SqlDatabase) InsertTxOuts(txs []*tsstypes.TxOutEntity) {
	query := "INSERT INTO tx_out (chain, hash_without_sig, in_chain, in_hash, bytes_without_sig, contract_hash) VALUES "
	query = query + getQueryQuestionMark(len(txs), 6)

	params := make([]interface{}, 0, len(txs)*6)

	for _, tx := range txs {
		params = append(params, tx.OutChain)
		params = append(params, tx.HashWithoutSig)
		params = append(params, tx.InChain)
		params = append(params, tx.InHash)
		params = append(params, tx.BytesWithoutSig)
		params = append(params, tx.ContractHash)
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
			ContractHash:    contractHash.String,
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
	log.Debug("Updating txout hash:", hash, " to status:", string(status), "chain", chain)
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
