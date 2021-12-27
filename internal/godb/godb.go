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
		rows.Scan(&user.Name, &user.Message)
		users = append(users, user)
	}

	fmt.Println(users)
	return users

}

func (i *Instance) updateUserMessage(ctx context.Context, name string, message string) {
	_, err := i.Db.Exec(ctx, "UPDATE users SET message=$1 WHERE name=$2;", message, name)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (i *Instance) getUserByName(ctx context.Context, name string) {
	user := &User{}
	err := i.Db.QueryRow(ctx, "SELECT name, message FROM users WHERE name=$1 LIMIT 1;", name).Scan(&user.Name, &user.Message)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("User by name: %v\nUser message: %v\n", user.Name, user.Message)
}
