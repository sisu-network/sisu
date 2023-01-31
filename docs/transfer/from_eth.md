To swap token from an EVM-based chain (ETH, binance, avalanche, etc), create a transferOut transaction and send it to Sisu vault contract deployed at the source blockchain.

The ABI of the vault contract could be [found](https://github.com/sisu-network/sisu/blob/master/contracts/eth/vault/Vault.abi) in Sisu contract code. Use transferOutNonEvm function call if you are transferring to a non-EVM chain.

Example of transferring from ETH based chain to Lisk could be found [here](https://github.com/sisu-network/transfer-examples/tree/master/ETH_to_Lisk).
