package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rubenv/sql-migrate"
	"log"
	"os"
	"net/http"
	"encoding/json"
	"strings"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

func main() {

	//Open sesame on the database
	getsettings()
	initialize_db()

	//Let's check for migrations and apply them
	migrate_db()

	//Routers
	http.HandleFunc("/auth/params", web_get_auth_params)
	http.HandleFunc("/auth/sign_in", web_sign_in)
	http.HandleFunc("/auth", web_register)
	//http.HandleFunc("/auth/change_pw",web_change_password)
	http.HandleFunc("/items/sync", sync_items)

	http.ListenAndServe(":8080", nil)
}

func web_get_auth_params(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if len(email) != 0 {
		row := db.QueryRow("SELECT pw_func,pw_alg,pw_cost,pw_key_size,pw_nonce,email,encrypted_password,uuid from users where email = $1", email)
		user := map_to_user(row, email)
		params := user.get_auth_params()
		js, err := json.Marshal(params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func web_sign_in(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are supported", 405)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var jsondata web_sign_in_struct
	err := decoder.Decode(&jsondata)
	if err != nil {
		panic(err)
	}
	email := jsondata.Email
	password := jsondata.Password
	if len(email) != 0 && len(password) != 0 {
		row := db.QueryRow("SELECT pw_func,pw_alg,pw_cost,pw_key_size,pw_nonce,email,encrypted_password,uuid from users where email = $1", email)
		user := map_to_user(row, email)
		token := user.sign_in(password)
		js, err := json.Marshal(sign_in_params{User: user, Token: token})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func web_register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are supported", 405)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var jsondata web_register_struct
	err := decoder.Decode(&jsondata)
	if err != nil {
		panic(err)
	}
	email := jsondata.Email
	password := jsondata.Password
	pw_alg := jsondata.Pw_alg
	pw_cost := jsondata.Pw_cost
	pw_func := jsondata.Pw_func
	pw_key_size := jsondata.Pw_key_size
	pw_nonce := jsondata.Pw_nonce

	user := User{Email: email, pw_alg: pw_alg, pw_cost: int32(pw_cost), pw_func: pw_func, pw_key_size: int32(pw_key_size), pw_nonce: pw_nonce}

	user = user.new_user(email, password)
	_, err = db.Exec("INSERT INTO USERS (email,encrypted_password,pw_alg,pw_cost,pw_func,pw_key_size,pw_nonce,uuid) values ($1,$2,$3,$4,$5,$6,$7,$8)", email, user.encrypted_password, user.pw_alg, user.pw_cost, user.pw_func, user.pw_key_size, user.pw_nonce, user.Uuid)
	if err != nil {
		panic(err)
	}

	token := user.sign_in(password)
	js, err := json.Marshal(sign_in_params{User: user, Token: token})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func migrate_db() {
	migrations := &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      "db/models",
	}
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		panic(err)
	}
	log.Printf("Applied %d migrations!\n", n)
}

func getsettings() {
	DB_HOST = os.Getenv("DB_HOST")
	DB_DATABASE = os.Getenv("DB_DATABASE")
	DB_USERNAME = os.Getenv("DB_USERNAME")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	SECRET_KEY_BASE = os.Getenv("SECRET_KEY_BASE")
	SALT_PSUEDO_NONCE = os.Getenv("SALT_PSUEDO_NONCE")
	log.Print("Settings read")
}

func initialize_db() {
	var err error
	db, err = sql.Open("postgres", "postgres://"+DB_USERNAME+":"+DB_PASSWORD+"@"+DB_HOST+"/"+DB_DATABASE+"?sslmode=disable")

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Print("Database connection successfull")

}

func sync_items(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var jsondata SyncRequest
	err := decoder.Decode(&jsondata)
	if err != nil {
		panic(err)
	}

	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		log.Print("Wrong Bearer")
	}

	token, err := jwt.Parse(authHeaderParts[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return SigningKey, nil
	})

	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

	}

}

func map_to_user(row *sql.Row, email string) User {
	user := User{Email: email}
	err := row.Scan(&user.pw_func, &user.pw_alg, &user.pw_cost, &user.pw_key_size, &user.pw_nonce, &user.Email, &user.encrypted_password, &user.Uuid)
	if err == sql.ErrNoRows {
		return User{Email: email}
	} else if err == nil {
		return user
	} else {
		panic(err)
	}
}

type web_register_struct struct {
	Email       string
	Password    string
	Pw_func     string
	Pw_alg      string
	Pw_key_size int32
	Pw_cost     int32
	Pw_salt     string
	Pw_nonce    string
}

type web_sign_in_struct struct {
	Email    string
	Password string
}

//SyncRequest - type for incoming sync request
type SyncRequest struct {
	Items       Items  `json:"items"`
	SyncToken   string `json:"sync_token"`
	CursorToken string `json:"cursor_token"`
	Limit       int    `json:"limit"`
}
