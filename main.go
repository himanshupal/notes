package main

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
)

type NullString string

type Tag NullString

func (s *NullString) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	strVal, ok := value.(string)
	if !ok {
		return errors.New("column is not a string")
	}
	*s = NullString(strVal)
	return nil
}

func (s NullString) Value() (driver.Value, error) {
	if len(s) == 0 { // if nil or empty string
		return nil, nil
	}
	return string(s), nil
}

type User struct {
	Id        string     `json:"_id,omitempty"`
	Username  string     `json:"username,omitempty"`
	FirstName NullString `json:"firstName,omitempty"`
	LastName  NullString `json:"lastName,omitempty"`
	Email     NullString `json:"email,omitempty"`
	Password  string     `json:"password,omitempty"`
	CreatedAt string     `json:"createdAt,omitempty"`
}

type Note struct {
	Id        string     `json:"_id,omitempty"`
	Tags      []Tag      `json:"tags,omitempty"`
	Title     string     `json:"title,omitempty"`
	Content   NullString `json:"content,omitempty"`
	Author    User       `json:"author,omitempty"`
	CreatedAt string     `json:"createdAt,omitempty"`
}

type Credentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func main() {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("PQ_HOST"), os.Getenv("PQ_PORT"), os.Getenv("PQ_USER"), os.Getenv("PQ_PASSWORD"), os.Getenv("PQ_DATABASE"))

	app := fiber.New(fiber.Config{
		AppName: "Notes",
	})

	app.Use(cors.New())
	app.Use(etag.New())
	app.Use(cache.New())
	app.Use(logger.New())
	app.Use(recover.New())

	app.Use(favicon.New(favicon.Config{
		File: "./assets/favicon.ico",
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: os.Getenv("COOKIE_KEY"),
	}))

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic("Couldn't connect to database: " + err.Error())
	}
	defer db.Close()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": fiber.StatusOK,
		})
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		token, err := Register(c, db)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "auth-token",
			HTTPOnly: true,
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 24 * 3),
		})

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"token": token,
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		data, err := Login(c, db)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "auth-token",
			HTTPOnly: true,
			Value:    fmt.Sprint(data["token"]),
			Expires:  time.Now().Add(time.Hour * 24 * 3),
		})

		return c.Status(fiber.StatusCreated).JSON(data)
	})

	// app.Get("/notes", func(c *fiber.Ctx) error {
	// 	res, err := getNotes(c, db)
	// })

	// app.Get("/notes/:id", func(c *fiber.Ctx) error {
	// 	res, err := getNote(c, db)
	// })

	// app.Post("/notes/:id", func(c *fiber.Ctx) error {
	// 	res, err := getNote(c, db)
	// })

	// app.Patch("/notes/:id", func(c *fiber.Ctx) error {
	// 	res, err := updateNote(c, db)
	// })

	// app.Delete("/notes/:id", func(c *fiber.Ctx) error {
	// 	res, err := deleteNote(c, db)
	// })

	// app.Get("/user", func(c *fiber.Ctx) error {
	// 	res, err := userInfo(c, db)
	// })

	// app.Patch("/user", func(c *fiber.Ctx) error {
	// 	res, err := updateUser(c, db)
	// })

	// app.Delete("/user", func(c *fiber.Ctx) error {
	//  res, err := deleteUser(c, db)
	// })

	if err := app.Listen(os.Getenv("HOST")); err != nil {
		panic("Couldn't start server: " + err.Error())
	}
}
