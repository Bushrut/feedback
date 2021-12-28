package godb

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Instance struct {
	Db *pgxpool.Pool
}

type User struct {
	Name    string
	Message string
}

func (i *Instance) Start() {
	fmt.Println("Progect godb started")
}

func (i *Instance) AddUser(ctx context.Context, name string, message string) {
	_, err := i.Db.Exec(ctx, "INSERT INTO users (created_at, name, message) VALUES ($1, $2, $3)", time.Now(), name, message)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Добавлен пользователь: %s\nСообщение: %s\n", name, message)
}

func (i *Instance) GetAllUsers(ctx context.Context) []User {
	var users []User
	rows, err := i.Db.Query(ctx, "SELECT name, message FROM users;")
	if err == pgx.ErrNoRows {
		fmt.Println("No rows")
	} else if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()

	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.Name, &user.Message)
		if err != nil {
			fmt.Println(err)
		}
		users = append(users, user)
	}

	fmt.Println(users)
	return users

}
