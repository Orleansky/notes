package api

import (
	"Anastasia/notes/pkg/storage"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Программный интерфейс сервера
type API struct {
	db     storage.Interface
	router *mux.Router
}

// Конструктор объекта API
func New(db storage.Interface) *API {
	api := API{
		db: db,
	}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

func (api *API) endpoints() {
	api.router.HandleFunc("/notes", api.getNotesHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/notes", api.createNoteHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/users", api.createUserHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/users", api.loginHandler).Methods(http.MethodGet, http.MethodOptions)
}

func (api *API) Router() *mux.Router {
	return api.router
}

// Получение всех заметок
func (api *API) getNotesHandler(w http.ResponseWriter, r *http.Request) {
	if storage.CurrentUserID == -1 {
		http.Error(w, "Для просмотра заметок войдите в систему", http.StatusNetworkAuthenticationRequired)
		log.Print("Не удалось просмотреть заметки, войдите в систему")
		return
	}
	posts, err := api.db.GetNotes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		log.Print("Не удалось просмотреть заметки, попробуйте еще раз")
		return
	}
	bytes, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print("Не удалось просмотреть заметки, попробуйте еще раз")
		return
	}

	log.Print("Получен доступ к вашим заметкам")
	w.Write(bytes)
}

// Добавление заметки
func (api *API) createNoteHandler(w http.ResponseWriter, r *http.Request) {
	if storage.CurrentUserID == -1 {
		http.Error(w, "Для добавления заметки войдите в систему", http.StatusNetworkAuthenticationRequired)
		log.Print("Не удалось добавить заметку, войдите в систему")
		return
	}
	var n storage.Note
	err := json.NewDecoder(r.Body).Decode(&n)
	if err != nil {
		http.Error(w, "Не удалось добавить заметку, попробуйте еще раз", http.StatusInternalServerError)
		log.Print("Не удалось добавить заметку, попробуйте еще раз")
		return
	}
	err = api.db.CreateNote(n)
	if err != nil {
		http.Error(w, "Не удалось добавить заметку, попробуйте еще раз", http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	log.Print("Заметка добавлена")
	w.WriteHeader(http.StatusOK)
}

func (api *API) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var u storage.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = api.db.CreateUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Print("Аккаунт создан")
	w.WriteHeader(http.StatusCreated)
}

func (api *API) loginHandler(w http.ResponseWriter, r *http.Request) {
	var u storage.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var ok bool

	ok, storage.CurrentUserID, err = api.db.LoginUser(u)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		log.Print("Вход не выполнен: неверный логин или пароль")
		return
	}

	log.Print("Вы вошли в систему")
	w.WriteHeader(http.StatusOK)
}
