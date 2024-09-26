package main

import (
	"Anastasia/notes/pkg/api"
	"Anastasia/notes/pkg/storage"
	"Anastasia/notes/pkg/storage/postgres"
	"log"
	"net/http"
)

type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// 	Создаём объект сервера.
	var srv server

	// Создаём строку подключения к PostgreSQL.
	connstr := "postgres://postgres:qwerty@localhost/notes?sslmode=disable"

	// Реляционная БД PostgreSQL.
	db, err := postgres.New(connstr)
	if err != nil {
		log.Fatal(err)
	}

	srv.db = db

	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	http.ListenAndServe(":8080", srv.api.Router())
}
