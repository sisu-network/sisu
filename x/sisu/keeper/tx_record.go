package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"
)

var (
	NotFound = errors.New("Not Found")
)

func GetTxRecordHash(msg sdk.Msg) ([]byte, string, error) {
	switch msg := msg.(type) {
	case *types.KeygenWithSigner:
		serialized, err := msg.Data.Marshal()
		if err != nil {
			return nil, "", err
		}
		return []byte(utils.KeccakHash32(string(serialized))), msg.Signer, nil

	case *types.KeygenResultWithSigner:
		hash := getKeygenKey(msg.Keygen.KeyType, int(msg.Keygen.Index))
		return hash, msg.Signer, nil

	case *types.TransfersMsg:
		serialized, err := msg.Data.Marshal()
		if err != nil {
			return nil, "", err
		}
		return []byte(utils.KeccakHash32(string(serialized))), msg.Signer, nil

	case *types.TxOutMsg:
		serialized, err := msg.Data.Marshal()
		if err != nil {
			return nil, "", err
		}
		return []byte(utils.KeccakHash32(string(serialized))), msg.Signer, nil

	case *types.TxOutResultMsg:
		serialized, err := msg.Data.Marshal()
		if err != nil {
			return nil, "", err
		}
		return []byte(utils.KeccakHash32(string(serialized))), msg.Signer, nil

	case *types.PauseContractMsg:
		serialized, err := msg.Data.Marshal()
		if err != nil {
			return nil, "", err
		}
		return []byte(utils.KeccakHash32(string(serialized))), msg.Signer, nil

	case *types.ResumeContractMsg:
		serialized, err := msg.Data.Marshal()
		if err != nil {
			return nil, "", err
		}
		return []byte(utils.KeccakHash32(string(serialized))), msg.Signer, nil

	case *types.ChangeOwnershipContractMsg:
		serialized, err := msg.Data.Marshal()
		if err != nil {
			return nil, "", err
		}
		return []byte(utils.KeccakHash32(string(serialized))), msg.Signer, nil

	case *types.LiquidityWithdrawFundMsg:
		serialized, err := msg.Data.Marshal()
		if err != nil {
			return nil, "", err
		}
		return []byte(utils.KeccakHash32(string(serialized))), msg.Signer, nil

	case *types.TransferFailureMsg:
		serialized, err := msg.Data.Marshal()
		if err != nil {
			return nil, "", err
		}
		return []byte(utils.KeccakHash32(string(serialized))), msg.Signer, nil

	case *types.TxInMsg:
		serialized, err := msg.Data.Marshal()
		if err != nil {
			return nil, "", err
		}
		return []byte(utils.KeccakHash32(string(serialized))), msg.Signer, nil
	}

	return nil, "", NotFound
}
