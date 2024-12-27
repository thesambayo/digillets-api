CREATE TABLE IF NOT EXISTS users (
  id bigserial PRIMARY KEY,
  public_id VARCHAR(50) UNIQUE,
  name text NOT NULL,
  email citext UNIQUE NOT NULL,
  password_hash bytea NOT NULL,
  activated bool NOT NULL DEFAULT false,
  created_at timestamp(0)
  with
    time zone NOT NULL DEFAULT NOW (),
    updated_at timestamp(0)
  with
    time zone NOT NULL DEFAULT NOW (),
    version integer NOT NULL DEFAULT 1
);
