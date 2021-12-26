CREATE TABLE IF NOT EXISTS tx_in(
  chain VARCHAR(256),
  hash VARCHAR(64),
  block_height BIGINT,
  serialized BLOB,

  PRIMARY KEY (chain, hash)
);
