CREATE TABLE IF NOT EXISTS keygen(
  chain VARCHAR(256),
  address VARCHAR(256),
  pubkey BLOB,
  status VARCHAR(64),

  PRIMARY KEY (chain)
);