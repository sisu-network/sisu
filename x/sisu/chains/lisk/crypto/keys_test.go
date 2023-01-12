package crypto

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	passphrase = "camera accident escape cricket frog pony record occur broken inhale waste swing"
	//address    = "lsk9hxtj8busjfugaxcg9zfuzdty7zyagcrsxvohk"
)

func TestGetPublicKeyFromSecret(t *testing.T) {
	publicKey, _ := hex.DecodeString("f0321539a45078365c1a65944d010876c0efe45c0446101dacced7a2f29aa289")
	val := GetPublicKeyFromSecret(passphrase)
	require.Equal(t, val, publicKey)
}

func TestGetPrivateKeyFromSecret(t *testing.T) {
	privateKey, _ := hex.DecodeString("ba20a2df297ff5db79764c7b4e778e00eaa81b665b551447fae4fdcd64c81b76f0321539a45078365c1a65944d010876c0efe45c0446101dacced7a2f29aa289")

	val := GetPrivateKeyFromSecret(passphrase)
	require.Equal(t, val, privateKey)
}
