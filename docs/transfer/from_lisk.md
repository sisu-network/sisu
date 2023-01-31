To swap Lisk tokens from Lisk testnet to another chain, you can send a transaction like a normal lisk transaction. The data field is where you define the information of destination chain, recipient and amount transfer

## Transfer Payload

The payload is a 64-based of the transfer protobuf. Lisk only allows max 64 characters in the payload and hence we have to use every field concisely.
```
message TransferData {
  required uint64 chainId = 1;
  required bytes recipient = 2;
  optional string token = 3;
  required uint64 amount = 4;
}
```
- Chain Id is the id of the destination chain.
- Recipient is the compacted string of the recipient address (each chain has a different way to compact the address).
- Token is an optional field reserved to indicate transfer token in the future For now Sisu only supports Liks token transfer and you can omit this field.
- Amount is the amount of Beddow (1 Lisk = 10^8 Beddow) to transfer.

## Compacted Address for Ethereum.
An Ethereum address has a hex format starting with `0x`. For example, `0xc275dc8be39f50d12f66b6a63629c39da5bae5bd`. The compacted format of EVM-based address is the decoded string of the hex substring after `0x`.

```
    recipient: Uint8Array.from(
      Buffer.from(recipientAddress.substring(2, recipientAddress.length), "hex")
    ),
```
The byte array takes much smaller buffer space than the original string.

## Examples
You can take a look at a JS example of how to swap Lisk token from [Lisk to ETH](https://github.com/sisu-network/transfer-examples/blob/master/Lisk_to_ETH/index.ts) chain.
