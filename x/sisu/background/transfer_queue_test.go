package background

// func mockForTransferQueue() (sdk.Context, ManagerContainer) {
// 	ctx := testmock.TestContext()
// 	k := testmock.KeeperTestGenesis(ctx)
// 	privateDb := keeper.NewPrivateDb(".", db.MemDBBackend)
// 	params := k.GetParams(ctx)
// 	params.TransferOutParams = []*types.TransferOutParams{
// 		{
// 			Chain:       "ganache2",
// 			MaxBatching: 1,
// 		},
// 	}
// 	k.SaveParams(ctx, params)

// 	txOutputProducer := &MockTxOutputProducer{}
// 	globalData := &components.MockGlobalData{
// 		GetReadOnlyContextFunc: func() sdk.Context {
// 			return ctx
// 		},
// 	}
// 	txSubmit := &components.MockTxSubmit{}
// 	mockAppKeys := &components.MockAppKeys{}
// 	valManagers := &components.MockValidatorManager{}

// 	mc := MockManagerContainer(ctx, k, txOutputProducer, globalData, txSubmit, privateDb, mockAppKeys,
// 		valManagers)

// 	return ctx, mc
// }

// func TestTransferQueue(t *testing.T) {
// 	t.Run("transfer_is_saved", func(t *testing.T) {
// 		ctx, mc := mockForTransferQueue()
// 		txOutProducer := mc.TxOutProducer().(*MockTxOutputProducer)
// 		appKeys := mc.AppKeys()

// 		txSubmit := mc.TxSubmit().(*components.MockTxSubmit)
// 		txSubmitCount := 0

// 		valManager := mc.ValidatorManager().(*components.MockValidatorManager)
// 		valManager.GetAssignedValidatorFunc = func(ctx sdk.Context, hash string) *types.Node {
// 			return &types.Node{
// 				AccAddress: appKeys.GetSignerAddress().String(),
// 			}
// 		}

// 		k := mc.Keeper()

// 		queue := NewTransferQueue(k, mc.TxOutProducer(), txSubmit,
// 			appKeys, mc.PrivateDb(), mc.ValidatorManager(),
// 			mc.GlobalData()).(*defaultTransferQueue)

// 		transfer := &types.TransferDetails{
// 			Id:          "ganache1__hash1",
// 			ToRecipient: "0x98Fa8Ab1dd59389138B286d0BeB26bfa4808EC80",
// 			Token:       "SISU",
// 			Amount:      utils.EthToWei.String(),
// 		}

// 		k.AddTransfers(ctx, []*types.TransferDetails{transfer})
// 		k.SetTransferQueue(ctx, "ganache2", []*types.TransferDetails{transfer})

// 		txOutProducer.GetTxOutsFunc = func(ctx sdk.Context, chain string,
// 			transfers []*types.TransferDetails) ([]*types.TxOutMsg, error) {
// 			ret := make([]*types.TxOutMsg, len(transfers))
// 			for i := range transfers {
// 				ret[i] = &types.TxOutMsg{
// 					Signer: appKeys.GetSignerAddress().String(),
// 					Data: &types.TxOut{
// 						Content: &types.TxOutContent{
// 							OutChain: "ganache2",
// 						},
// 					},
// 				}
// 			}
// 			return ret, nil
// 		}

// 		txSubmit.SubmitMessageAsyncFunc = func(msg sdk.Msg) error {
// 			txSubmitCount++
// 			return nil
// 		}

// 		queue.processBatch(ctx)
// 		require.Equal(t, 1, txSubmitCount)
// 	})
// }
