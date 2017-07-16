package main

import (
	"time"
)

type User struct {
	Uuid               string    `json:"uuid"`
	Email              string    `json:"email"`
	pw_func            string
	pw_alg             string
	pw_cost            int32
	pw_key_size        int32
	pw_nonce           string
	encrypted_password string
	created_at         time.Time
	updated_at         time.Time
}

type pw_params struct {
	Pw_func     string `json:"pw_func"`
	Pw_alg      string `json:"pw_alg"`
	Pw_salt     string `json:"pw_salt"`
	Pw_cost     int32  `json:"pw_cost"`
	Pw_key_size int32  `json:"pw_key_size"`
}

type sign_in_params struct {
	User  User `json:"user"`
	Token string `json:"token"`
}

type Item struct {
	Uuid         string    `json:"uuid"`
	User_uuid    string    `json:"user_uuid"`
	Content      string    `json:"content"`
	Content_type string    `json:"content_type"`
	Enc_item_key string    `json:"enc_item_key"`
	Auth_hash    string    `json:"auth_hash"`
	Deleted      bool      `json:"deleted"`
	Created_at   time.Time `json:"created_at"`
	Updated_at   time.Time `json:"updated_at"`
}

//Items - is an items slice
type Items []Item

//SyncRequest - type for incoming sync request
type SyncRequest struct {
	Items       Items  `json:"items"`
	SyncToken   string `json:"sync_token"`
	CursorToken string `json:"cursor_token"`
	Limit       int    `json:"limit"`
}

type unsaved struct {
	Item
	error
}

//SyncResponse - type for response
type SyncResponse struct {
	Retrieved   Items     `json:"retrieved_items"`
	Saved       Items     `json:"saved_items"`
	Unsaved     []unsaved `json:"unsaved"`
	SyncToken   string    `json:"sync_token"`
	CursorToken string    `json:"cursor_token,omitempty"`
}
