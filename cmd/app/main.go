package main

import (
	"context"
	"feedback/internal/godb"
	"feedback/pkg/helpers/pg"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

var tpl = template.Must(template.ParseFiles("web/index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func addFeedback(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	nickname := r.Form.Get("nickname")
	msg := r.Form.Get("message")

	fmt.Println(nickname, msg)
	tpl.Execute(w, nil)

}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/addfeedback", addFeedback)
	http.ListenAndServe(":8080", mux)

	cfg := &pg.Config{
		Host:     "localhost",
		Username: "db_user",
		Password: "pwd123",
		Port:     "54320",
		DbName:   "db_test",
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

	ins := &godb.Instance{Db: c}
	ins.Start()

}
