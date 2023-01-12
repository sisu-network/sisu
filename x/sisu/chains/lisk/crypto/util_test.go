package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSHA256Hash(t *testing.T) {
	val := GetSHA256Hash(passphrase)
	passphraseHash := "ba20a2df297ff5db79764c7b4e778e00eaa81b665b551447fae4fdcd64c81b76"
	require.Equal(t, hex.EncodeToString(val[:]), passphraseHash)

}

func TestGetFirstEightBytesReversed(t *testing.T) {
	defaultBytes := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	bytesReversed := []byte{7, 6, 5, 4, 3, 2, 1, 0}

	successVal := GetFirstEightBytesReversed(defaultBytes)
	require.Equal(t, successVal, bytesReversed)

	failedVal := GetFirstEightBytesReversed(nil)
	require.Equal(t, failedVal, []byte(nil))
}

func TestGetBigNumberStringFromBytes(t *testing.T) {
	defaultBytes := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	bytesBigNumberString := "18591708106338011145"

	val := GetBigNumberStringFromBytes(defaultBytes)
	require.Equal(t, val, bytesBigNumberString)
}
