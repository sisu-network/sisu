package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func compileContract(name string) {
	// solc --abi name.sol --overwrite -o name
	_, err := exec.Command(
		"solc",
		"--abi",
		fmt.Sprintf("%s.sol", name),
		"--overwrite",
		"-o",
		fmt.Sprintf("%s", strings.ToLower(name)),
	).Output()
	if err != nil {
		panic(err)
	}

	// solc --bin name.sol --overwrite -o name
	_, err = exec.Command(
		"solc",
		"--bin",
		fmt.Sprintf("%s.sol", name),
		"--overwrite",
		"-o",
		fmt.Sprintf("%s", strings.ToLower(name)),
	).Output()
	if err != nil {
		panic(err)
	}

	// abigen --bin=name/name.bin --abi=name/name.abi --pkg=name --out=name/name.go
	_, err = exec.Command(
		"abigen",
		fmt.Sprintf("--bin=%s/%s.bin", name, name),
		fmt.Sprintf("--abi=%s/%s.abi", name, name),
		fmt.Sprintf("--pkg=%s", name),
		fmt.Sprintf("--out=%s/%s.go", strings.ToLower(name), strings.ToLower(name)),
	).Output()
	if err != nil {
		panic(err)
	}
}

// Generates contracts
func main() {
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".sol" {
			return nil
		}

		parts := strings.Split(path, ".")
		compileContract(parts[0])

		return nil
	})
}
