package sisu

// func TestProcessor_OnTxIns(t *testing.T) {
// 	t.Parallel()

// 	t.Run("empty_tx", func(t *testing.T) {
// 		t.Parallel()

// 		processor := &Processor{}
// 		require.NoError(t, processor.OnTxIns(&eyesTypes.Txs{}))
// 	})

// 	t.Run("success_to_our_key", func(t *testing.T) {
// 		t.Parallel()

// 		ctrl := gomock.NewController(t)
// 		t.Cleanup(func() {
// 			ctrl.Finish()
// 		})

// 		mockTxSubmit := mock.NewMockTxSubmit(ctrl)
// 		mockTxSubmit.EXPECT().SubmitMessageAsync(gomock.Any()).Return(nil).Times(1)

// 		observedChain := "eth"
// 		toAddress := utils.RandomHeximalString(64)
// 		fromAddres := utils.RandomHeximalString(64)

// 		mockPublicDb := mocktss.NewMockStorage(ctrl)
// 		mockPublicDb.EXPECT().IsKeygenAddress(libchain.KEY_TYPE_ECDSA, fromAddres).Return(false).Times(1)

// 		priv := ed25519.GenPrivKey()
// 		addr := sdk.AccAddress(priv.PubKey().Address())
// 		appKeysMock := mock.NewMockAppKeys(ctrl)
// 		appKeysMock.EXPECT().GetSignerAddress().Return(addr).MinTimes(1)

// 		txs := &eyesTypes.Txs{
// 			Chain: observedChain,
// 			Block: int64(utils.RandomNaturalNumber(1000)),
// 			Arr: []*eyesTypes.Tx{{
// 				Hash:       utils.RandomHeximalString(64),
// 				Serialized: []byte{},
// 				To:         toAddress,
// 				From:       fromAddres,
// 			}},
// 		}

// 		// Init processor with mocks
// 		processor := &Processor{
// 			publicDb: mockPublicDb,
// 			appKeys:  appKeysMock,
// 			txSubmit: mockTxSubmit,
// 		}

// 		err := processor.OnTxIns(txs)
// 		// <-done

// 		require.NoError(t, err)
// 	})
// }
