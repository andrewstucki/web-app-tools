CREATE TABLE IF NOT EXISTS memberships (
  namespace_id uuid NOT NULL,
  user_id uuid NOT NULL,
  role varchar(50) NOT NULL,
  PRIMARY KEY (namespace_id, user_id)
);
