package storage

type Note struct {
	ID       int
	Title    string
	Content  string
	UserID   int
	Username string
}

type User struct {
	ID       int
	Name     string
	Password string
}

var CurrentUserID = -1

type Interface interface {
	GetNotes() ([]Note, error)         // получение списка заметок
	CreateNote(Note) error             // добавление заметки
	CreateUser(User) error             // добавление пользователя
	LoginUser(User) (bool, int, error) // аутентификация пользователя
}
