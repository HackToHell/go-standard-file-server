package main

import (
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"crypto/sha256"
	"encoding/hex"
	"crypto/sha1"
	"github.com/nu7hatch/gouuid"
	"time"
	"fmt"
	"log"
	"database/sql"
)

func (user User) sign_in(password string) string {
	err := bcrypt.CompareHashAndPassword([]byte(user.encrypted_password), []byte(password))
	if err != nil {
		panic(err)
	} else {
		hasher := sha256.New()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_uuid": user.Uuid,
			"pw_hash":   hex.EncodeToString(hasher.Sum([]byte(user.encrypted_password))),
		})
		tokenizedstring, err := token.SignedString([]byte(SECRET_KEY_BASE))
		if err != nil {
			panic(err)
		}
		return tokenizedstring
	}

}

func (user User) get_auth_params() pw_params {
	hasher := sha1.New()
	if user.pw_nonce != "" && user.pw_cost != 0 && user.pw_alg != "" && user.pw_key_size != 0 && user.pw_func != "" {
		return pw_params{
			Pw_func:     user.pw_func,
			Pw_alg:      user.pw_alg,
			Pw_salt:     hex.EncodeToString(hasher.Sum([]byte(user.Email + "@" + string(user.pw_nonce)))),
			Pw_cost:     user.pw_cost,
			Pw_key_size: user.pw_key_size,
		}
	} else {
		return pw_params{
			Pw_func:     "pbkdf2",
			Pw_alg:      "sha512",
			Pw_salt:     hex.EncodeToString(hasher.Sum([]byte(user.Email + "@" + SALT_PSUEDO_NONCE))),
			Pw_cost:     5000,
			Pw_key_size: 512,
		}
	}
}

func (user User) new_user(email string, password string) User {
	var err error
	u, err := uuid.NewV4()
	user.Uuid = u.String()
	user.Email = email
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	user.encrypted_password = string(pass)
	if err != nil {
		panic(err)
	}
	return user
	//user.pw_func = "pbkdf2"
	//user.pw_alg = "sha512"
	//user.pw_cost = 5000
	//user.pw_nonce = SALT_PSUEDO_NONCE
	//user.pw_key_size = 512
}

func (items Items) sync(sync_token string, cursor_token string, limit int){

}

func (this *Item) create() error {
	if this.Uuid == "" {
		tmp,err := uuid.NewV4()
		if err != nil {
			panic(err)
		}
		this.Uuid = tmp.String()

	}
	this.Created_at = time.Now()
	this.Updated_at = time.Now()
	_, err :=  db.Query("INSERT INTO `items` (`uuid`, `user_uuid`, content,  content_type, enc_item_key, auth_hash, deleted, created_at, updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)", this.Uuid, this.User_uuid, this.Content, this.Content_type, this.Enc_item_key, this.Auth_hash, this.Deleted, this.Created_at, this.Updated_at)
	return err
}

func (this *Item) update() error {
	this.Updated_at = time.Now()
	_,err := db.Query("UPDATE `items` SET `content`=$1, `enc_item_key`=$2, `auth_hash`=$3, `deleted`=$4, `updated_at`=$5 WHERE `uuid`=$6 AND `user_uuid`=$6", this.Content, this.Enc_item_key, this.Auth_hash, this.Deleted, this.Updated_at, this.Uuid, this.User_uuid)
	return err
}

func (this *Item) delete() error {
	if this.Uuid == "" {
		return fmt.Errorf("Trying to delete unexisting item")
	}
	this.Content = ""
	this.Enc_item_key = ""
	this.Auth_hash = ""
	this.Updated_at = time.Now()

	_,err := db.Query("UPDATE `items` SET `content`='', `enc_item_key`='', `auth_hash`='',`deleted`=1, `updated_at`=? WHERE `uuid`=? AND `user_uuid`=?", this.Updated_at, this.Uuid, this.User_uuid)
	return err
}

func (this Item) copy() (Item, error) {
	_tmp,err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	this.Uuid = _tmp.String()
	this.Updated_at = time.Now()
	err = this.create()
	if err != nil {
		log.Print(err)
		return Item{}, err
	}
	return this, nil
}

func (this Item) Exists() bool {
	if this.Uuid == "" {
		return false
	}
	var uuid string
	err := db.QueryRow("SELECT `uuid` FROM `items` WHERE `uuid`=$1", this.Uuid).Scan(&uuid)

	if err != nil {
		log.Print(err)
		return false
	}
	return uuid != ""
}

func (this *Item) LoadByUUID(uuid string) bool {
	row := db.QueryRow("SELECT * FROM `items` WHERE `uuid`=$1", uuid)
	err := map_to_item(row,this)

	if err != nil {
		log.Print(err)
		return false
	}

	return true
}

func map_to_item(row *sql.Row, item *Item) error {
	err := row.Scan(&item.Uuid,&item.Content_type,&item.Enc_item_key,&item.Auth_hash,&item.User_uuid,&item.Created_at,&item.Updated_at,&item.Deleted)
		return err
}