package database

var CreateSchema = `
CREATE TABLE curiosity.user (
     id           SERIAL NOT NULL,
     google_id    CHARACTER VARYING(64),
     username     CHARACTER VARYING(50) NOT NULL,
     password     CHARACTER VARYING(100) NOT NULL,
     email        CHARACTER VARYING(50) NOT NULL,
     hashed_email CHARACTER VARYING(64) NOT NULL,
     is_active    BOOLEAN NOT NULL,
     private_key  CHARACTER VARYING(4000),
     public_key   CHARACTER VARYING(1000),
     PRIMARY KEY (id),
     UNIQUE (username)
);
`

var DropSchema = `
DROP TABLE IF EXISTS curiosity.user CASCADE;
`
