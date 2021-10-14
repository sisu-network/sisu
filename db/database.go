package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	tsstypes "github.com/sisu-network/sisu/x/tss/types"
)

type Database interface {
	Init() error
	Close() error

	// Contracts
	InsertConctracts(contracts []*tsstypes.ContractEntity)
	GetPendingDeployContracts(chain string) []*tsstypes.ContractEntity
	UpdateContractsState(contracts []*tsstypes.ContractEntity, state string)
	UpdateContractDeployTx(chain, id string, txHash string)

	// Txout
	InsertTxOuts(txs []*tsstypes.TxOutEntity)
	UpdateTxOut(chain string, hashWithoutSig string, bz []byte, sig []byte)
	IsContractDeployTx(chain string, hashWithoutSig string) bool
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
		utils.LogError("Failed to connect to DB. Err =", err)
		return err
	}

	err = d.DoMigration()
	if err != nil {
		utils.LogError("Cannot do migration. Err =", err)
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
	utils.LogInfo("Db is connected successfully")
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

	utils.LogInfo("Doing sql migration...")

	m.Log = &dbLogger{}
	m.Up()

	return nil
}

func (d *SqlDatabase) InsertConctracts(contracts []*tsstypes.ContractEntity) {
	query := "INSERT INTO contract (chain, hash, name) VALUES "
	query = query + getQueryQuestionMark(len(contracts), 3)

	params := make([]interface{}, 0, 3*len(contracts))
	for _, contract := range contracts {
		params = append(params, contract.Chain)
		params = append(params, contract.Hash)
		params = append(params, contract.Name)
	}

	_, err := d.db.Exec(query, params...)
	if err != nil {
		utils.LogError("failed to insert contract into db, err = ", err)
	}
}

func (d *SqlDatabase) GetPendingDeployContracts(chain string) []*tsstypes.ContractEntity {
	query := "SELECT chain, hash, name, state FROM contract WHERE chain=?"
	params := []interface{}{chain}
	result := make([]*tsstypes.ContractEntity, 0)

	rows, err := d.db.Query(query, params...)
	if err != nil {
		return result
	}

	for rows.Next() {
		var chain, hash, name, state sql.NullString
		if err := rows.Scan(&chain, &hash, &name, &state); err != nil {
			utils.LogError("cannot scan row, err =", err)
			continue
		}

		if state.String == "" {
			result = append(result, &tsstypes.ContractEntity{
				Chain: chain.String,
				Hash:  hash.String,
				Name:  name.String,
				State: state.String,
			})
		}
	}

	return result
}

func (d *SqlDatabase) UpdateContractsState(contracts []*tsstypes.ContractEntity, state string) {
	for _, contract := range contracts {
		query := "UPDATE contract SET state = ? WHERE chain = ? AND hash = ?"
		params := []interface{}{state, contract.Chain, contract.Hash}

		_, err := d.db.Exec(query, params...)
		if err != nil {
			utils.LogError("failed to update contract state, err =", err, ". len(contracts) =", len(contracts))
		}
	}
}

func (d *SqlDatabase) UpdateContractDeployTx(chain, id string, txHash string) {
	query := "UPDATE contract SET tx_hash = ? WHERE chain = ? AND id = ?"
	params := []interface{}{txHash, chain, id}

	_, err := d.db.Exec(query, params...)
	if err != nil {
		utils.LogError("failed to update contract deploy tx, err =", err)
	}
}

func (d *SqlDatabase) InsertTxOuts(txs []*tsstypes.TxOutEntity) {
	query := "INSERT INTO tx_out (chain, hash_without_sig, in_chain, in_hash, bytes, contract_hash) VALUES "
	query = query + getQueryQuestionMark(len(txs), 6)

	params := make([]interface{}, 0, len(txs)*6)

	for _, tx := range txs {
		params = append(params, tx.OutChain)
		params = append(params, tx.HashWithoutSig)
		params = append(params, tx.InChain)
		params = append(params, tx.InHash)
		params = append(params, tx.Outbytes)
		params = append(params, tx.ContractHash)
	}

	_, err := d.db.Exec(query, params...)
	if err != nil {
		utils.LogError("failed to insert txout into table, err = ", err)
	}
}

func (d *SqlDatabase) UpdateTxOut(chain string, hashWithoutSig string, bz []byte, sig []byte) {
	query := "UPDATE tx_out SET bytes = ?, signature = ? WHERE chain = ? AND hash_without_sig = ?"
	params := []interface{}{
		bz,
		sig,
		chain,
		hashWithoutSig,
	}

	_, err := d.db.Exec(query, params...)
	if err != nil {
		utils.LogError("failed to insert txout with chain and hashWoSig", chain, hashWithoutSig, ", err =", err)
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

	if rows.Next() {
		var hash sql.NullString
		if err := rows.Scan(&hash); err != nil {
			return false
		}

		return hash.String != ""
	}

	return false
}

func (d *SqlDatabase) Close() error {
	return d.db.Close()
}
