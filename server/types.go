package main

import (
    "database/sql"
    "net/http"
)

type Storer interface {
    Init() error

    CreateDatabaseINE(string) error

    CreateUserTableINE() error
    CreateUser(ReqCreateUser) (User, error)
    DeleteUser(id uint) error

    CreateTodoTableINE() error
    CreateTodo(ReqCreateTodo) (Todo, error)
    DeleteTodo(id uint) error

    GetUserTodos(id uint) (*sql.Rows, error)
}

type Server struct {
    DB   Storer
    Addr string
    Mux  *http.ServeMux
}

func NewServer(addr string, db *sql.DB, mux *http.ServeMux) Server {
    return Server {
        Mux: mux,
        Addr: addr,
        DB: &PostgresStorage{db},
    }
}

type User struct {
    Id        uint   `json:"id"`
    Username  string `json:"username"`
    FirstName string `json:"firstName"`
    LastName  string `json:"lastName"`
}

type Todo struct {
    Id     uint   `json:"id"`
    Title  string `json:"title"`
    UserId uint   `json:"userId"`
}

type ReqCreateUser struct {
    Username  string `json:"username"`
    FirstName string `json:"firstName"`
    LastName  string `json:"lastName"`
    Password  string `json:"password"`
}

type ReqCreateTodo struct {
    Title  string `json:"title"`
    UserId string `json:"userId"`
}

type ReqGetUserTodos struct {
    UserId uint   `json:"id"`
    Title  string `json:"title"`
}


