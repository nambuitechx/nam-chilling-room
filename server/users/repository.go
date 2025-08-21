package users

import (
	"database/sql"
	"errors"
)

type UserRepository struct {
	DB	*sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) selectUsers(username string, limit int, offset int) ([]User, error) {
	var users = []User{}

	rows, err := r.DB.Query(
		"SELECT id, username, password, created_at, updated_at FROM users WHERE username LIKE $1 LIMIT $2 OFFSET $3",
		"%" + username + "%",
		limit,
		offset,
	)
	
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) selectUserByID(id string) (*User, error) {
	var user User

	row := r.DB.QueryRow(
		"SELECT id, username, password, created_at, updated_at FROM users WHERE id = $1",
		id,
	)

	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) selectUserByUsername(username string) (*User, error) {
	var user User

	row := r.DB.QueryRow(
		"SELECT id, username, password, created_at, updated_at FROM users WHERE username = $1",
		username,
	)

	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) insertUser(user *User) (*User, error) {
	var insertedUser User

	err := r.DB.QueryRow(
		`INSERT INTO users(id, username, password)
		VALUES($1, $2, $3)
		RETURNING id, username, password, created_at, updated_at
		`,
		user.ID,
		user.Username,
		user.Password,
	).Scan(&insertedUser.ID, &insertedUser.Username, &insertedUser.Password, &insertedUser.CreatedAt, &insertedUser.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &insertedUser, nil
}
