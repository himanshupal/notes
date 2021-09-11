package main

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// Registers a new user account & login
func Register(c *fiber.Ctx, db *sql.DB) (string, error) {
	credentials := new(Credentials)

	if err := c.BodyParser(credentials); err != nil {
		return "", err
	}

	if credentials.Username == "" {
		return "", fiber.NewError(fiber.StatusBadRequest, "Username is required!")
	}

	if len(credentials.Username) < 3 {
		return "", fiber.NewError(fiber.StatusBadRequest, "Username must be 3 or more characters!")
	}

	if credentials.Password == "" {
		return "", fiber.NewError(fiber.StatusBadRequest, "Password is required!")
	}

	if len(credentials.Password) < 11 {
		return "", fiber.NewError(fiber.StatusBadRequest, "Password must be 11 or more characters!")
	}

	var userExists bool

	db.QueryRow(`SELECT COUNT("username") FROM "users" WHERE "username"=$1`, credentials.Username).Scan(&userExists)
	if userExists {
		return "", fiber.NewError(fiber.StatusConflict, "Username already exists!")
	}

	var id string

	hashedPassword, err := GenerateHash(credentials.Password, DefaultParams)
	if err != nil {
		return "", err
	}

	if err := db.QueryRow(`INSERT INTO "users" ("username", "password") VALUES ($1, $2) RETURNING id`, credentials.Username, hashedPassword).Scan(&id); err != nil {
		return "", err
	}

	authToken, err := GenerateToken(id, 24*3)
	if err != nil {
		return "", err
	}

	return authToken, nil
}
