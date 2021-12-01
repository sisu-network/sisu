CREATE TABLE IF NOT EXISTS keygen(
  key_type VARCHAR(64),
  address VARCHAR(256),
  pubkey BLOB,
  status VARCHAR(64),
  start_block BIGINT,

  PRIMARY KEY (key_type)
);
