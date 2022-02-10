package main

import "github.com/jmoiron/sqlx"

type Repo struct {
	conn *sqlx.DB
}

func NewRepo(conn *sqlx.DB) *Repo {
	return &Repo{
		conn: conn,
	}
}

type User struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func (r *Repo) CreateUser(name string) (user User, err error) {
	err = r.conn.Get(&user, "INSERT INTO users(name) VALUES ($1) RETURNING *", name)
	return
}

func (r *Repo) GetAllUser() (users []User, err error) {
	err = r.conn.Select(&users, "SELECT * FROM users")
	return
}

func RunMigrations(conn *sqlx.DB) error {
	_, err := conn.Exec(`
	CREATE TABLE IF NOT EXISTS users 
		(
		    id   SERIAL PRIMARY KEY, 
		    name TEXT NOT NULL UNIQUE
		)
	`)
	return err
}
