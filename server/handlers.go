package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func Register(s *Server) {

	s.Mux.HandleFunc("POST /make-user", s.HandleCreateUser)
	s.Mux.HandleFunc("DELETE /delete-user", s.HandleDeleteUser)

	s.Mux.HandleFunc("POST /make-todo", s.HandleCreateTodo)
	s.Mux.HandleFunc("DELETE /delete-todo", s.HandleDeleteTodo)

    s.Mux.HandleFunc("GET /user-todos", s.HandleGetUserTodos)
    s.Mux.HandleFunc("GET /test-ping", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Hello World")
    })
}

func (s *Server) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var data ReqCreateUser
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        log.Println("[ERROR]", err)
        w.WriteHeader(http.StatusBadRequest)
        return
	}

	err := s.DB.CreateUser(data)
	if err != nil {
        log.Println("[ERROR]", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
	}
    log.Println("[INFO]", "POST /make-user")

    w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")

    id, err := strconv.Atoi(idStr)
    if err != nil || id < 1 {
        log.Println("[ERROR]", err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    err = s.DB.DeleteUser(uint(id))
    if err != nil {
        log.Println("[ERROR]", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    log.Println("[INFO]", "DELETE /delete-user")
    w.WriteHeader(http.StatusOK)
}


func (s *Server) HandleCreateTodo(w http.ResponseWriter, r *http.Request) {
	var data ReqCreateTodo
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        log.Println("[ERROR]", err)
        w.WriteHeader(http.StatusBadRequest)
        return
	}

	err := s.DB.CreateTodo(data)
	if err != nil {
        log.Println("[ERROR]", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
	}
    log.Println("[INFO]", "POST /make-todo")
    w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleDeleteTodo(w http.ResponseWriter, r *http.Request) {
    var data struct { Id uint `json:"id"` }

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        log.Println("[ERROR]", err)
        w.WriteHeader(http.StatusBadRequest)
        return
	}

    err := s.DB.DeleteTodo(uint(data.Id))
    if err != nil {
        log.Println("[ERROR]", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    log.Println("[INFO]", "DELETE /delete-user")
    w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleGetUserTodos(w http.ResponseWriter, r *http.Request) {
    var data struct { Id uint `json:"userId"` }

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        log.Println("[ERROR]", err)
        w.WriteHeader(http.StatusBadRequest)
        return
	}
    if data.Id <= 0 {
        log.Println("[ERROR]", "request with invalid id")
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    rows, err := s.DB.GetUserTodos(data.Id)
    if err != nil {
        log.Println("[ERROR]", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    sendData := []ReqGetUserTodos{}
    for rows.Next() {
        var id uint
        var title string
        err := rows.Scan(&id, &title)
        if err != nil {
            log.Println("[ERROR]", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        sendData = append(sendData, ReqGetUserTodos{ UserId: id, Title: title })
    }

    json, err := json.Marshal(sendData)
    if err != nil {
        log.Println("[ERROR]", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    log.Println("[INFO]", "GET /user-todos")
    w.WriteHeader(http.StatusOK)
    w.Write(json)
}




