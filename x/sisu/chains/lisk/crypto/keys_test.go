package crypto

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	passphrase = "camera accident escape cricket frog pony record occur broken inhale waste swing"
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

func TestAddressFromPublicKey(t *testing.T) {
	publicKey, _ := hex.DecodeString("aae33db460ee8e037ea87e5f8f0de34e138ef118c4b7113b53cdf3a1f6618e90")
	val := GetAddressFromPublicKey(publicKey)
	require.Equal(t, val, "750c899a72aae6911a003645acf7fc7459c29d35")
}

func TestAddressToLisk32(t *testing.T) {
	testValues := []map[string]string{
		{"address": "0db493f33fd4cee095d14050137299ec05f3d8e1", "lisk32": "lskxed4njagdtn7xm7y3x3xbjkahuvgnenxm98zau"},
		{"address": "1a0047f16007d13b911ff93eaf10dff0723a7dfd", "lisk32": "lskc3zp8j5zzg3twp3ggpg6fpbgfxackg8hvao6gf"},
		{"address": "3d322f86d1c4a575ff71f5fabb3e2deb2fc3d0ef", "lisk32": "lsknkavgxey2rrw5gsyfwh5e8y9howjnkn8w4epcj"},
		{"address": "19b46edbba9c2cd214bd709c86ea960705117783", "lisk32": "lskcbdbhehdtue9pmmh7v739dkezjvyvhjcunb2v8"},
		{"address": "3a6eb418cd4628b1b894996da6152452ddf38afb", "lisk32": "lskn4w53bb932k5c7pktmedom4p657gnvwstae5m9"},
	}
	for _, testVal := range testValues {
		address, _ := hex.DecodeString(testVal["address"])
		val := AddressToLisk32(address)
		require.Equal(t, val, testVal["lisk32"])
	}
}
