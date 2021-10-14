CREATE TABLE chain_key(
  chain VARCHAR(256),
  address VARCHAR(256),
  pubkey BLOB,

  PRIMARY KEY (chain, address)
);
