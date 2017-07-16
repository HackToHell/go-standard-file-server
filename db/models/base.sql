
-- +migrate Up
CREATE TABLE items (
  uuid uuid NOT NULL PRIMARY KEY ,
content TEXT,
content_type TEXT DEFAULT NULL,
enc_item_key TEXT,
auth_hash varchar(255) DEFAULT NULL,
user_uuid uuid DEFAULT NULL,
created_at TIMESTAMP(6) NOT NULL,
updated_at TIMESTAMP(6) NOT NULL,
deleted smallint DEFAULT '0'
);


CREATE index  ON items(updated_at);
CREATE index  ON items(user_uuid);
CREATE index  ON items(user_uuid,content_type);


CREATE TABLE users (
uuid UUID NOT NULL PRIMARY KEY ,
email varchar(255) DEFAULT NULL,
pw_func varchar(255) DEFAULT NULL,
pw_alg varchar(255) DEFAULT NULL,
pw_cost bigint DEFAULT NULL,
pw_key_size bigint DEFAULT NULL,
pw_nonce varchar(255) DEFAULT NULL,
encrypted_password varchar(255) NOT NULL DEFAULT '',
created_at TIMESTAMP DEFAULT NULL,
updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX ON users(email);