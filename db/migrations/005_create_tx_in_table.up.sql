CREATE TABLE IF NOT EXISTS tx_in(
  chain VARCHAR(256),
  hash VARCHAR(64),
  serialized BLOB,

  PRIMARY KEY (chain, hash)
);
