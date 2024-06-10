package main

import (
    "database/sql"
    "errors"
    "fmt"
    "log"

    _ "github.com/lib/pq"
)

const (
    dbPort = "5432"
)

func Connect() (*sql.DB, error) {
    connStr := fmt.Sprintf("user=postgres dbname=test port=%v password=sid sslmode=disable", dbPort)

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, errors.New("[ERROR] connecting to database: " + err.Error())
    }

    err = db.Ping()
    if err != nil {
        return nil, errors.New("[ERROR] while pinging the database: " + err.Error())
    }

    return db, nil
}

func Disconnect(db *sql.DB) error {
    err := db.Close()
    if err != nil {
        log.Println("[ERROR]", err.Error())
    }
    return err
}

type PostgresStorage struct {
    DB *sql.DB
}

func (d *PostgresStorage) Init() error {
    err := d.CreateUserTableINE()
    if err != nil { return err }

    err = d.CreateTodoTableINE()
    if err != nil { return err }

    return nil
}

func (d *PostgresStorage) CreateDatabaseINE(dbName string) error {
    query :=
    `CREATE DATABAdE IF NOT EXISTS $1`

    _, err := d.DB.Query(query, dbName)
    return errors.New("[ERROR]" + "error creating a database " + dbName + err.Error())
}

func (d *PostgresStorage) CreateUserTableINE() error {
    query := 
    `CREATE TABLE IF NOT EXISTS "user" (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    password TEXT NOT NULL)`

    if _, err := d.DB.Query(query); err != nil {
        return errors.New("error creating \"user\" table: " + err.Error())
    }
    return nil
}

func (d *PostgresStorage) CreateTodoTableINE() error {
    query := 
    `CREATE TABLE IF NOT EXISTS todo (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    user_id INT NOT NULL,
    CONSTRAINT fk_todo_userid 
    FOREIGN KEY(user_id) 
    REFERENCES "user"(id))`

    if _, err := d.DB.Query(query); err != nil {
        return errors.New("error creating todo table: " + err.Error())
    }
    return nil
}

func (d *PostgresStorage) CreateUser(r ReqCreateUser) (User, error) {
    query := 
    `INSERT INTO "user" (username, first_name, last_name, password)
    VALUES ($1, $2, $3, $4) RETURNING id, username, first_name, last_name`

    var user User 
    err := d.DB.QueryRow(query, r.Username, r.FirstName, r.LastName, r.Password).Scan(
        &user.Id, 
        &user.Username, 
        &user.FirstName, 
        &user.LastName,
    )
    return user, err
}

func (d *PostgresStorage) DeleteUser(id uint) error {
    query := `DELETE FROM "user" WHERE id=$1`

    _, err := d.DB.Query(query, id)
    return err
}


func (d *PostgresStorage) CreateTodo(r ReqCreateTodo) (Todo, error) {
    query := 
    `INSERT INTO todo (title, user_id)
    VALUES ($1, $2) RETURNING id, title, user_id`

    var todo Todo
    err := d.DB.QueryRow(query, r.Title, r.UserId).Scan(
        &todo.Id, 
        &todo.Title, 
        &todo.UserId, 
    )
    return todo, err
}

func (d *PostgresStorage) DeleteTodo(id uint) error {
    query := `DELETE FROM "todo" WHERE id=$1`

    _, err := d.DB.Query(query, id)
    return err
}

func (d *PostgresStorage) GetUserTodos(id uint) (*sql.Rows, error){
    log.Println("[DBGINFO]", "id:", id)
    query := 
    `SELECT todo.id, todo.title 
    FROM todo 
    JOIN "user"
    ON todo.user_id = "user".id
    WHERE todo.user_id = $1`

    rows, err := d.DB.Query(query, id)
    return rows, err
}


