package gen

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"text/template"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	econfig "github.com/sisu-network/deyes/config"
	heartconfig "github.com/sisu-network/dheart/core/config"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
	tmos "github.com/tendermint/tendermint/libs/os"
)

type DeyesConfiguration struct {
	Chains        []econfig.Chain
	SisuServerUrl string

	// sql
	Sql    SqlConfig
	LogDNA log.LogDNAConfig `toml:"log_dna"`
}

func getPeerIds(n int, pubKeys []cryptotypes.PubKey) ([]string, error) {
	ids := make([]string, n)

	for i := 0; i < n; i++ {
		p2pPubKey, err := crypto.UnmarshalEd25519PublicKey(pubKeys[i].Bytes())
		if err != nil {
			panic(err)
		}

		id, err := peer.IDFromPublicKey(p2pPubKey)
		if err != nil {
			panic(err)
		}

		ids[i] = id.String()
	}

	return ids, nil
}

func writeDeyesConfig(deyesConfig DeyesConfiguration, dir string) {
	eyesToml := `db_host = "{{ .Sql.Host }}"
db_port = {{ .Sql.Port }}
db_username = "{{ .Sql.Username }}"
db_password = "{{ .Sql.Password }}"
db_schema = "{{ .Sql.Schema }}"

server_port = 31001
sisu_server_url = "{{ .SisuServerUrl }}"

[log_dna]
secret = "{{ .LogDNA.Secret }}"
app_name = "{{ .LogDNA.AppName }}"
host_name = "{{ .LogDNA.HostName }}"
flush_interval = "{{ .LogDNA.FlushInterval }}"
max_buffer_len = {{ .LogDNA.MaxBufferLen }}
log_local = {{ .LogDNA.LogLocal }}

[chains]{{ range $k, $chain := .Chains }}
[chains.{{ $chain.Chain }}]
  chain = "{{ $chain.Chain }}"
  block_time = {{ $chain.BlockTime }}
  adjust_time = {{ $chain.AdjustTime }}
  starting_block = 0
  rpcs = [{{ range $j, $rpc := $chain.Rpcs }} "{{ $rpc }}" {{end}}]
  rpc_secret = "{{ $chain.RpcSecret }}"{{ end }}
`

	tmpl := template.New("eyesToml")

	configTemplate, err := tmpl.Parse(eyesToml)
	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer
	err = configTemplate.Execute(&buffer, deyesConfig)

	tmos.MustWriteFile(filepath.Join(dir, "deyes.toml"), buffer.Bytes(), 0644)
}

func writeHeartConfig(dir string, heartConfig heartconfig.HeartConfig) {

	heartToml := `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

home-dir = "/root/"
use-on-memory = {{ .UseOnMemory }}
shortcut-preparams = {{ .ShortcutPreparams }}
sisu-server-url = "{{ .SisuServerUrl }}"
port = 5678

[log_dna]
secret = "{{ .LogDNA.Secret }}"
app_name = "{{ .LogDNA.AppName }}"
host_name = "{{ .LogDNA.HostName }}"
flush_interval = "{{ .LogDNA.FlushInterval }}"
max_buffer_len = {{ .LogDNA.MaxBufferLen }}
log_local = {{ .LogDNA.LogLocal }}


###############################################################################
###                        Database Configuration                           ###
###############################################################################
[db]
  host = "{{ .Db.Host }}"
  port = {{ .Db.Port }}
  username = "{{ .Db.Username }}"
  password = "{{ .Db.Password }}"
  schema = "{{ .Db.Schema }}"
[connection]
  host = "{{ .Connection.Host }}"
  port = {{ .Connection.Port }}
  rendezvous = "{{ .Connection.Rendezvous }}"
{{ range $k, $v := .Connection.Peers }}
[[connection.peers]]
	address = {{ $v.Address }}
	pubkey = "{{ $v.PubKey }}"
	pubkey_type = "{{ $v.PubKeyType }}"
{{ end }}
`

	tmpl := template.New("heartToml")

	configTemplate, err := tmpl.Parse(heartToml)
	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer
	err = configTemplate.Execute(&buffer, heartConfig)

	tmos.MustWriteFile(filepath.Join(dir, "dheart.toml"), buffer.Bytes(), 0644)
}

func getChains(file string) []*types.Chain {
	chains := []*types.Chain{}

	dat, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(dat, &chains); err != nil {
		panic(err)
	}

	return chains
}

func getTokens(file string) []*types.Token {
	tokens := []*types.Token{}

	dat, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(dat, &tokens); err != nil {
		panic(err)
	}

	return tokens
}

func getLiquidity(file string) []*types.Liquidity {
	liquids := []*types.Liquidity{}
	bz, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(bz, &liquids); err != nil {
		panic(err)
	}

	return liquids
}
