-- this table acts as a singleton for the site settings

CREATE TABLE IF NOT EXISTS site_settings (
  initialized boolean NOT NULL DEFAULT FALSE
);

INSERT INTO site_settings (initialized) VALUES (FALSE);
