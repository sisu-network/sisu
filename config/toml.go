package config

import (
	"bytes"
	"text/template"

	tmos "github.com/tendermint/tendermint/libs/os"
)

const defaultConfigTemplate = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

mode = "{{ .Mode }}"

###############################################################################
###                     Siu Main App Connfiguration                         ###
###############################################################################
[sisu]
chain-id = "{{ .Sisu.ChainId}}"
keyring-backend = "{{ .Sisu.KeyringBackend }}"
api-host = "{{ .Sisu.ApiHost }}"
api-port = {{ .Sisu.ApiPort }}

###############################################################################
###                         Siu TSS Connfiguration                          ###
###############################################################################
[tss]
dheart-host = "{{ .Tss.DheartHost }}"
dheart-port = {{ .Tss.DheartPort }}
deyes-url = "{{ .Tss.DeyesUrl }}"
[tss.supported-chains] {{ range $k, $v := .Tss.SupportedChains }}
	[tss.supported-chains.{{ $v.Id }}]
		id = "{{ $v.Id }}"
		token = "{{ $v.Token }}"
{{ end }}
`

var configTemplate *template.Template

func init() {
	var err error

	tmpl := template.New("appConfigFileTemplate")

	if configTemplate, err = tmpl.Parse(defaultConfigTemplate); err != nil {
		panic(err)
	}
}

func WriteConfigFile(configFilePath string, config *Config) {
	var buffer bytes.Buffer

	if err := configTemplate.Execute(&buffer, config); err != nil {
		panic(err)
	}

	tmos.MustWriteFile(configFilePath, buffer.Bytes(), 0644)
}
