CREATE TABLE IF NOT EXISTS sessions(
  token TEXT,
  data BYTEA NOT NULL,
  expiry TIMESTAMPTZ NOT NULL,
  CONSTRAINT sessions_pkey PRIMARY KEY (token)
);
