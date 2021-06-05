package config

import "io/ioutil"

func writeDefaultTss(filePath string) error {
	var defaultValue = `enable = false
host = "localhost"
port = 5678
[supported_chains]
`
	err := ioutil.WriteFile(filePath, []byte(defaultValue), 0644)
	if err != nil {
		return err
	}

	return nil
}
