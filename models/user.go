package models

import (
	"errors"
	"fmt"
	"stock-dashboard/db"
	"stock-dashboard/utils"
	"strings"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role,omitempty"`
}

type UserUpdate struct {
	ID       string `json:"id"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
}

func (u *User) ValidateCredentials() error {
	query := `
		SELECT id, password, role FROM users WHERE email = $1
	`

	row := db.DB.QueryRow(query, strings.ToLower(u.Email))

	var retrivedPassword string
	err := row.Scan(&u.ID, &retrivedPassword, &u.Role)
	if err != nil {
		return err
	}

	isValidPassword := utils.CheckPasswordHash(u.Password, retrivedPassword)

	if !isValidPassword {
		return errors.New("invalid email or password")
	}
	return nil
}

func (u *User) Save() error {
	if u.Role == "" {
		u.Role = "staff"
	}

	hashedPassword, err := utils.HashPassword(u.Password)

	if err != nil {
		return err
	}

	query := `INSERT INTO users(email,password,role) VALUES($1,$2,$3) RETURNING id`

	err = db.DB.QueryRow(query, strings.ToLower(u.Email), hashedPassword, u.Role).Scan(&u.ID)

	if err != nil {
		return err
	}

	return err
}

func GetAllStaff() ([]User, error) {
	query := `
		SELECT id, email, role FROM users WHERE role = 'staff' ORDER BY email
	`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u *User) Get() error {
	query := `
		SELECT id, email, role FROM users WHERE id = $1
	`

	row := db.DB.QueryRow(query, u.ID)
	err := row.Scan(&u.ID, &u.Email, &u.Role)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUpdate) Update() error {
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", u.ID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user not found")
	}

	query := "UPDATE users SET "
	args := []interface{}{}
	argIndex := 1

	if u.Email != "" {
		query += "email = $" + fmt.Sprintf("%d", argIndex) + ", "
		args = append(args, u.Email)
		argIndex++
	}

	if u.Password != "" {
		hashedPassword, err := utils.HashPassword(u.Password)
		if err != nil {
			return err
		}
		query += "password = $" + fmt.Sprintf("%d", argIndex) + ", "
		args = append(args, hashedPassword)
		argIndex++
	}

	if u.Role != "" {
		query += "role = $" + fmt.Sprintf("%d", argIndex) + ", "
		args = append(args, u.Role)
		argIndex++
	}

	query = query[:len(query)-2]
	query += " WHERE id = $" + fmt.Sprintf("%d", argIndex)
	args = append(args, u.ID)

	_, err = db.DB.Exec(query, args...)
	return err
}

func (u *User) Delete() error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := db.DB.Exec(query, u.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
