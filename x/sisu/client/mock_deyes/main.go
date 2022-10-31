package main

import (
	"context"
	"encoding/json"
	"math/big"
	"os"

	"github.com/mr-tron/base58"
	eyessolanatypes "github.com/sisu-network/deyes/chains/solana/types"
	solanatypes "github.com/sisu-network/sisu/x/sisu/chains/solana/types"

	solanago "github.com/gagliardetto/solana-go"
	eyestypes "github.com/sisu-network/deyes/types"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sisu-network/lib/log"
)

func getSolanaPrivateKey() solanago.PrivateKey {
	key, err := solanago.PrivateKeyFromSolanaKeygenFile(os.Getenv("KEY_PATH"))
	if err != nil {
		panic(err)
	}

	return key
}

func getSolanaTransferoutTx() *eyessolanatypes.Transaction {
	ix := &solanatypes.TransferOutInstruction{
		Instruction: byte(solanatypes.TranserOut),
		Data: solanatypes.TransferOutData{
			Amount:       *big.NewInt(9000),
			TokenAddress: "BfFFZs7e6n62rMKqYR9pLtQnVw3rPFDTosyKWB8xfgzs",
			ChainId:      189985,
			Recipient:    "0x8095f5b69F2970f38DC6eBD2682ed71E4939f988",
		},
	}

	bz, err := ix.Serialize()
	if err != nil {
		panic(err)
	}

	outerTx := &eyessolanatypes.Transaction{
		TransactionInner: &eyessolanatypes.TransactionInner{
			Signatures: []string{"Gcof1VG8san3BTyAknKXzqZtGLbEbiCVqnH1umc1EFFJVafChskwjGkJrUGD1geRYaLDRR4RW7C7mrrLijmTnJ3"},
			Message: &eyessolanatypes.TransactionMessage{
				AccountKeys: []string{
					"GWP9AoY6ZvUqLzm4fS5jqSJAJ8rnrMf4d1kiU1wSXwED",
					"EuFCgxwUQMFoC8N1iFazZsn1C66wyDx4dTP4RSpMJscw",
					"TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
					"AUzpPvijCvMiX6CTSbiby2nTU2sXjbEQWYTfRDjoH969",
					"8kBCKTsqi1FpCgUiigJCLa5PGyyyeXETxYAiSRnXRArX",
					"7nAXTnyYCT4QCixT2uo1d2Tgc9wGYnP1i6KfD6VEbKK9",
				},
				Instructions: []eyessolanatypes.Instruction{
					{
						ProgramIdIndex: 0,
						Data:           base58.Encode(bz),
					},
				},
			},
		},
	}

	return outerTx
}

func TestPostSolanaTx() {
	client, err := rpc.DialContext(context.Background(), "http://0.0.0.0:25456")
	if err != nil {
		panic(err)
	}

	outerTx := getSolanaTransferoutTx()
	bz, err := json.Marshal(outerTx)

	txs := &eyestypes.Txs{
		Chain:     "solana-devnet",
		Block:     172485419,
		BlockHash: "",
		Arr: []*eyestypes.Tx{
			{
				Hash:       "Gcof1VG8san3BTyAknKXzqZtGLbEbiCVqnH1umc1EFFJVafChskwjGkJrUGD1geRYaLDRR4RW7C7mrrLijmTnJ3",
				Serialized: bz,
				Success:    true,
			},
		},
	}

	var result string
	err = client.CallContext(context.Background(), &result, "tss_postObservedTxs", txs)
	if err != nil {
		panic(err)
	}
	log.Verbose("Done broadcasting!")
}

func main() {
	TestPostSolanaTx()
}
