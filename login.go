package main

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx, db *sql.DB) (fiber.Map, error) {
	user := new(User)
	credentials := new(Credentials)

	if err := c.BodyParser(credentials); err != nil {
		return nil, err
	}

	if credentials.Username == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Username is required!")
	}

	if len(credentials.Username) < 3 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Username must be 3 or more characters!")
	}

	if credentials.Password == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Password is required!")
	}

	if len(credentials.Password) < 11 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Password must be 11 or more characters!")
	}

	if err := db.QueryRow(`SELECT "id", "username", "first_name", "password" from "users" WHERE "username"=$1`, credentials.Username).Scan(&user.Id, &user.Username, &user.FirstName, &user.Password); err != nil {
		return nil, err
	}

	ok, _, err := VerifyHash(credentials.Password, user.Password)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fiber.NewError(fiber.StatusNotFound, "Incorrect password!")
	}

	token, err := GenerateToken(user.Id, 24*3)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return fiber.Map{
		"token":       token,
		"userDetails": user,
	}, nil
}
