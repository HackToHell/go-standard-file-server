package main

import "database/sql"

var DB_HOST string
var DB_DATABASE string
var DB_USERNAME string
var DB_PASSWORD string
var db *sql.DB
var SECRET_KEY_BASE string
var SALT_PSUEDO_NONCE string
