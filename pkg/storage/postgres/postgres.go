package postgres

import (
	"Anastasia/notes/pkg/speller"
	"Anastasia/notes/pkg/storage"
	"context"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Store struct {
	db *pgxpool.Pool
}

func New(connstr string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}

	s := Store{
		db: db,
	}

	return &s, nil
}

func (s *Store) CreateNote(note storage.Note) error {
	arrTitle := strings.Split(note.Title, " ")
	arrContent := strings.Split(note.Content, " ")
	text := []string{}

	text = append(text, arrTitle[:]...)
	text = append(text, arrContent[:]...)

	text, err := speller.CheckText(text)

	if err != nil {
		return err
	}

	arrTitle = text[:len(arrTitle)]
	arrContent = text[len(arrContent)+1:]

	title := strings.Join(arrTitle, " ")
	content := strings.Join(arrContent, " ")

	_, err = s.db.Exec(context.Background(), `
		INSERT INTO notes(user_id, title, content)
		VALUES
		($1, $2, $3);
	`, storage.CurrentUserID, title, content)

	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetNotes() ([]storage.Note, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT
		notes.id,
		notes.title,
		notes.content,
		users.id AS user_id,
		users.name
		FROM notes
		JOIN users ON notes.user_id = users.id
		WHERE users.id = $1
	`, storage.CurrentUserID)

	if err != nil {
		return nil, err
	}

	var notes []storage.Note

	for rows.Next() {
		var n storage.Note
		err = rows.Scan(
			&n.ID,
			&n.Title,
			&n.Content,
			&n.UserID,
			&n.Username,
		)

		if err != nil {
			return nil, err
		}

		notes = append(notes, n)
	}
	return notes, rows.Err()
}

func (s *Store) CreateUser(user storage.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(context.Background(), `
		INSERT INTO users (name, password)
		VALUES ($1, $2)
	`, user.Name, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) LoginUser(user storage.User) (bool, int, error) {
	var hashedPassword string

	row := s.db.QueryRow(context.Background(), `
		SELECT password
		FROM users
		WHERE name = $1
	`, user.Name)
	err := row.Scan(&hashedPassword)
	if err != nil {
		return false, -1, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		return false, -1, nil
	}

	if err != nil {
		return false, -1, err
	}

	var id int
	row = s.db.QueryRow(context.Background(), `
		SELECT id
		FROM users
		WHERE name = $1
	`, user.Name)
	err = row.Scan(&id)
	if err != nil {
		return false, -1, err
	}

	return true, id, nil
}
