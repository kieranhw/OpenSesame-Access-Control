CREATE TABLE IF NOT EXISTS entries (
  entry_id    SERIAL        PRIMARY KEY,
  name        VARCHAR(255)  NOT NULL,
  description TEXT,
  created_at  TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ   NOT NULL DEFAULT now()
);