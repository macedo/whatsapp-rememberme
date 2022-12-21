CREATE TABLE IF NOT EXISTS users(
  id UUID DEFAULT uuid_generate_v4(),
  username CITEXT NOT NULL UNIQUE,
  encrypted_password VARCHAR NOT NULL,
  CONSTRAINT users_pkey PRIMARY KEY (id)
);
