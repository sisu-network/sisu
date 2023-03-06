package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"
)

func compileContract(name, contract string) {
	// solc --abi name.sol --overwrite -o name
	output, err := exec.Command(
		"solc",
		"--abi",
		fmt.Sprintf("%s.sol", name),
		"--overwrite",
		"-o",
		fmt.Sprintf("%s", name),
	).Output()
	if err != nil {
		panic(string(output))
	}

	// solc --bin name.sol --overwrite -o name
	output, err = exec.Command(
		"solc",
		"--bin", fmt.Sprintf("%s.sol", name),
		"--overwrite",
		"-o", name,
	).Output()
	if err != nil {
		panic(string(output))
	}

	// abigen --bin=name/contract.bin --abi=name/contract.abi --pkg=name --out=name/name.go
	output, err = exec.Command(
		"abigen",
		fmt.Sprintf("--bin=%s/%s.bin", name, contract),
		fmt.Sprintf("--abi=%s/%s.abi", name, contract),
		fmt.Sprintf("--pkg=%s", name),
		fmt.Sprintf("--out=%s/%s.go", name, strings.ToLower(name)),
	).Output()
	if err != nil {
		panic(string(output))
	}
}

// Generates contracts
func main() {
	name := flag.String("name", "", "Name of the solidity file (without the extension).")
	contract := flag.String("contract", "", "Name of contract in solidity file.")

	flag.Parse()
	compileContract(*name, *contract)
}
