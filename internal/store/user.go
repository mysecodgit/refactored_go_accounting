package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Password string `json:"-"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) GetAll(ctx context.Context) ([]User, error) {
	query := `
		SELECT id, name, username, phone
		FROM users
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Username,
			&u.Phone,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

/*
GET BY ID
*/
func (s *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, name, username, phone
		FROM users
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var u User
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Name,
		&u.Username,
		&u.Phone,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (name, username, phone, password)
		VALUES (?, ?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.Username,
		user.Phone,
		user.Password,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (s *UserStore) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE users
		SET name = ?,
		    username = ?,
		    phone = ?
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.Username,
		user.Phone,
		user.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *UserStore) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM users
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
