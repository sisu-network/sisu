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
keyring-passphrase = "{{ .Sisu.KeyringPassphrase }}"
api-host = "{{ .Sisu.ApiHost }}"
api-port = {{ .Sisu.ApiPort }}
rpc-port = {{ .Sisu.RpcPort }}
internal-api-port = {{ .Sisu.InternalApiPort }}

###############################################################################
###                     Siu Main LogDNA Connfiguration                      ###
###############################################################################
[log_dna]
secret = "{{ .LogDNA.Secret }}"
app_name = "{{ .LogDNA.AppName }}"
host_name = "{{ .LogDNA.HostName }}"
flush_interval = "{{ .LogDNA.FlushInterval }}"
max_buffer_len = {{ .LogDNA.MaxBufferLen }}
log_local = {{ .LogDNA.LogLocal }}

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
