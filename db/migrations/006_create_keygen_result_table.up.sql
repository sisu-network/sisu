CREATE TABLE IF NOT EXISTS keygen_result(
  key_type VARCHAR(64),
  keygen_index INT,
  result int,

  PRIMARY KEY (key_type, keygen_index)
);
