CREATE TABLE IF NOT EXISTS contract(
  chain VARCHAR(256),
  hash VARCHAR(256),
  byte_code BLOB,
  name VARCHAR(256),
  address VARCHAR(256),
  status VARCHAR(256),

  PRIMARY KEY (chain, hash)
);
