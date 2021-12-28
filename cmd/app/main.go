package main

import (
	"context"
	"feedback/internal/godb"
	"feedback/pkg/helpers/pg"
	"fmt"
	"net/http"
	"os"
	"text/template"
)

var tpl = template.Must(template.ParseFiles("web/index.html"))
var user = &godb.User{}
var db = godb.Instance{}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.Execute(w, nil)
	if err != nil {
		fmt.Println(err)
	}

}

func addFeedback(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}

	nickname := r.Form.Get("nickname")
	msg := r.Form.Get("message")
	user.Name = nickname
	user.Message = msg
	db.AddUser(context.Background(), user.Name, user.Message)

	http.Redirect(w, r, "/", http.StatusFound)

}

func main() {

	cfg := &pg.Config{
		Host:     "localhost",
		Username: "db_user",
		Password: "pwd123",
		Port:     "54320",
		DbName:   "db_feedback",
		Timeout:  5,
	}

	poolConfig, err := pg.NewPoolConfig(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Pool config error: %v\n", err)
		os.Exit(1)
	}

	poolConfig.MaxConns = 5

	c, err := pg.NewConnection(poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connet to database failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Connection OK")

	_, err = c.Exec(context.Background(), ";")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ping failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Ping OK")

	db = godb.Instance{Db: c}
	db.Start()

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/addfeedback", addFeedback)
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err)
	}

}
