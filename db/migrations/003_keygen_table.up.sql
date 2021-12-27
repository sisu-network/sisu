CREATE TABLE IF NOT EXISTS keygen(
  key_type VARCHAR(64),
  keygen_index INT,
  address VARCHAR(256),
  pubkey BLOB,
  start_block BIGINT,

  PRIMARY KEY (key_type, keygen_index)
);
