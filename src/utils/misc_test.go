package utils

import (
	"fmt"
	"testing"
)

func TestKeccakHash32(t *testing.T) {
	s := "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"getName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newName\",\"type\":\"string\"}],\"name\":\"setName\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	hash, err := KeccakHash32(s)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Hash = ", hash)
}
