CREATE TABLE tx_out(
  chain VARCHAR(256),
  hash_without_sig VARCHAR(64),

  hash_with_sig VARCHAR(64),
  in_hash VARCHAR(64),
  in_chain VARCHAR(256),
  bytes BLOB,
  signature BLOB,
  contract_hash VARCHAR(64),

  PRIMARY KEY (chain, hash_without_sig)
);
