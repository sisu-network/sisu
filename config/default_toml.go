package config

import "os"

func writeDefaultTss(filePath string) error {
	var defaultValue = `enable = false
dheart_host = "localhost"
dheart_port = 5678
[supported_chains]
`
	err := os.WriteFile(filePath, []byte(defaultValue), 0644)
	if err != nil {
		return err
	}

	return nil
}
