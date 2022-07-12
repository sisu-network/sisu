package cardano

import (
	"fmt"
	"testing"
)

func Test_wordToByteString(t *testing.T) {
	ret := wordToByteString("WRAP_ADA")
	fmt.Println(ret)
}
