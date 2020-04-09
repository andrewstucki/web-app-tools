CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  email varchar(50) UNIQUE NOT NULL,
  google_id varchar(50) UNIQUE NOT NULL
);
