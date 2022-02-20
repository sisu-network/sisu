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
	"github.com/sisu-network/sisu/x/sisu/types"
	tmos "github.com/tendermint/tendermint/libs/os"
)

type HeartConfiguration struct {
	PeerString        string
	SisuServerUrl     string
	UseOnMemory       string
	ShortcutPreparams bool
	Sql               SqlConfig
}

type DeyesConfiguration struct {
	Chains        []ChainConfig
	SisuServerUrl string

	// sql
	Sql SqlConfig
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

[chains]{{ range $k, $chain := .Chains }}
[chains.{{ $chain.Id }}]
  chain = "{{ $chain.Id }}"
  block_time = 1000
  starting_block = 0
  rpc_url = "{{ $chain.Rpc }}"{{ end }}
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

func writeHeartConfig(index int, dir string, peerString string, useOnMemory string, shortcutPreparams bool, sisuUrl string, sqlConfig SqlConfig) {
	heartConfig := HeartConfiguration{
		PeerString:        peerString,
		SisuServerUrl:     sisuUrl,
		Sql:               sqlConfig,
		UseOnMemory:       useOnMemory,
		ShortcutPreparams: shortcutPreparams,
	}

	heartToml := `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

home-dir = "/root/"
use-on-memory = {{ .UseOnMemory }}
shortcut-preparams = {{ .ShortcutPreparams }}
sisu-server-url = "{{ .SisuServerUrl }}"
port = 5678

###############################################################################
###                        Database Configuration                           ###
###############################################################################
[db]
  host = "{{ .Sql.Host }}"
  port = {{ .Sql.Port }}
  username = "{{ .Sql.Username }}"
  password = "{{ .Sql.Password }}"
  schema = "{{ .Sql.Schema }}"
[connection]
  host = "0.0.0.0"
  port = 28300
  rendezvous = "rendezvous"
  peers = [{{ .PeerString }}]
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
