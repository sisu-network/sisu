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

func writeDeyesConfig(deyesConfig econfig.Deyes, dir string) {
	eyesToml := `db_host = "{{ .DbHost }}"
db_port = {{ .DbPort }}
db_username = "{{ .DbUsername }}"
db_password = "{{ .DbPassword }}"
db_schema = "{{ .DbSchema }}"

server_port = 31001
sisu_server_url = "{{ .SisuServerUrl }}"

use_external_rpcs_info = {{ .UseExternalRpcsInfo }}

[log_dna]
secret = "{{ .LogDNA.Secret }}"
app_name = "{{ .LogDNA.AppName }}"
host_name = "{{ .LogDNA.HostName }}"
flush_interval = "{{ .LogDNA.FlushInterval }}"
max_buffer_len = {{ .LogDNA.MaxBufferLen }}
log_local = {{ .LogDNA.LogLocal }}

[price_providers]{{ range $name, $provider := .PriceProviders }}
[price_providers.{{ $name }}]
  url = "{{ $provider.Url }}"
  secrets = "{{ $provider.Secrets }}"{{ end }}

[tokens]{{ range $name, $token := .Tokens }}
[tokens.{{ $name }}]
  symbol = "{{ $token.Symbol }}"
  coin_cap_name = "{{ $token.CoincapName }}"
  coin_gecko_name = "{{ $token.CoinGeckoName }}"{{ end }}

[chains]{{ range $k, $chain := .Chains }}
[chains.{{ $chain.Chain }}]
  chain = "{{ $chain.Chain }}"
  block_time = {{ $chain.BlockTime }}
  adjust_time = {{ $chain.AdjustTime }}
  starting_block = 0
  rpcs = [{{ range $j, $rpc := $chain.Rpcs }}"{{ $rpc }}", {{end}}]
  wss = [{{ range $j, $ws := $chain.Wss }}"{{ $ws }}", {{end}}]
  use_eip_1559 = {{ $chain.UseEip1559 }}
  rpc_secret = "{{ $chain.RpcSecret }}"
  client_type = "{{ $chain.ClientType }}"{{ if $chain.SyncDB.Host }}
  [chains.{{ $chain.Chain }}.sync_db]
    host = "{{ $chain.SyncDB.Host }}"
    port = {{ $chain.SyncDB.Port }}
    user = "{{ $chain.SyncDB.User }}"
    password = "{{ $chain.SyncDB.Password }}"
    db_name = "{{ $chain.SyncDB.DbName }}"
    submit_url = "{{ $chain.SyncDB.SubmitURL }}"{{end}}
  solana_bridge_program_id="{{ $chain.SolanaBridgeProgramId }}"{{ end }}
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

func updateOracleSecret(deyesCfg econfig.Deyes) econfig.Deyes {
	for name, provider := range deyesCfg.PriceProviders {
		if len(provider.Secrets) != 0 {
			continue
		}

		switch name {
		case "coin_cap":
			provider.Secrets = os.Getenv("COIN_CAP_SECRET")
		case "coin_market_cap":
			provider.Secrets = os.Getenv("COIN_MARKET_CAP_SECRET")
		}

		deyesCfg.PriceProviders[name] = provider
	}

	return deyesCfg
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

func getVaults(file string) []*types.Vault {
	vaults := []*types.Vault{}
	bz, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(bz, &vaults); err != nil {
		panic(err)
	}

	return vaults
}
