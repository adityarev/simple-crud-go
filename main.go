package main

import (
	"database/sql"
	"github.com/gadp22/crema"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func GetUser(conditions map[string]string) (*sql.Rows, error) {
	q := crema.GetGenericSelectQuery("users", conditions, "id")

	return crema.ExecuteQuery(q.QueryString)
}

func PostUser(tx *sql.Tx, values map[string]string) *sql.Row {
	q := crema.GetGenericInsertQuery("users", values)

	return crema.ExecuteQueryRow(tx, q.QueryString)
}

func PutUser(tx *sql.Tx, values map[string]string) (sql.Result, error) {
	q := crema.GetGenericUpdateQuery("users", values, "id")

	return crema.ExecuteNonQuery(q.QueryString)
}

func DeleteUser(tx *sql.Tx, conditions map[string]string) (sql.Result, error) {
	q := crema.GetGenericDeleteQuery("users", conditions, "id")

	return crema.ExecuteNonQuery(q.QueryString)
}

func generateToken(w http.ResponseWriter, r *http.Request) {
	//give you own signing key (default signing key -> "crema")
	crema.SetSigningKey("halogaizz")

	//token expiration in minutes
	var expiration time.Duration = 360

	//any information you want to embed into the token
	test := make(map[string]string)
	test["data"] = "dafuq"

	tokenString, _ := crema.GenerateJWT(test, expiration)

	w.Write([]byte(tokenString))
}

func validateToken(w http.ResponseWriter, r *http.Request) {
	crema.SetSigningKey("halogaizz")

	vars := mux.Vars(r)
	tokenString := vars["token"]

	err := crema.ValidateJWT(tokenString)

	valid := "token " + tokenString + " is valid"

	if err != nil {
		valid = "token " + tokenString + " is invalid"
	}

	w.Write([]byte(valid))
}

func main() {
	server := crema.InitServer()

	server.AddRoutes(http.MethodGet, "/hello", Hello)

	server.AddRoutes(http.MethodGet, "/users", crema.MakeGenericGetHandler(GetUser))
	server.AddRoutes(http.MethodGet, "/users/{id}", crema.MakeGenericGetHandler(GetUser))
	server.AddRoutes(http.MethodPost, "/users", crema.MakeGenericPostHandler(PostUser))
	server.AddRoutes(http.MethodPut, "/users/{id}", crema.MakeGenericPutHandler(PutUser))
	server.AddRoutes(http.MethodDelete, "/users/{id}", crema.MakeGenericDeleteHandler(DeleteUser))

	server.AddRoutes(http.MethodGet, "/generateToken", generateToken)
	server.AddRoutes(http.MethodGet, "/validateToken/{token}", validateToken)

	crema.LogPrintf("[MAIN] Server is running, listening to port 8001 ....")
	log.Fatal(http.ListenAndServe(":8001", server.Router))
}
