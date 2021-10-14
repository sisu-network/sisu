CREATE TABLE contract(
  chain VARCHAR(256),
  hash VARCHAR(256),
  name VARCHAR(256),
  address VARCHAR(256),
  state VARCHAR(256),
  tx_deploy_hash VARCHAR(64),
  PRIMARY KEY (chain, hash)
);
