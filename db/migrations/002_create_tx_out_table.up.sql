CREATE TABLE IF NOT EXISTS tx_out(
  in_chain VARCHAR(64),
  in_hash VARCHAR(64),
  out_chain VARCHAR(64),
  out_hash VARCHAR(64),

  bytes_without_sig BLOB,
  status VARCHAR(64),
  hash_with_sig VARCHAR(64),
  signature BLOB,
  contract_hash VARCHAR(64),

  PRIMARY KEY (in_chain, out_chain, out_hash)
);
